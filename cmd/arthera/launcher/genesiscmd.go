package launcher

import (
	"compress/gzip"
	"errors"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"io"
	"math"
	"os"
	"path"
	"strconv"

	"github.com/artheranet/arthera-node/genesis"
	"github.com/artheranet/arthera-node/genesis/genesisstore"
	"github.com/artheranet/arthera-node/genesis/genesisstore/fileshash"
	"github.com/artheranet/arthera-node/gossip/evmstore"
	"github.com/artheranet/arthera-node/internal/inter/ibr"
	"github.com/artheranet/arthera-node/internal/inter/ier"
	"github.com/artheranet/arthera-node/utils/devnullfile"
	"github.com/artheranet/arthera-node/utils/iodb"
	"github.com/artheranet/lachesis/common/bigendian"
	"github.com/artheranet/lachesis/hash"
	"github.com/artheranet/lachesis/inter/idx"
	"github.com/artheranet/lachesis/kvdb"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

type dropableFile struct {
	io.ReadWriteSeeker
	io.Closer
	path string
}

func (f dropableFile) Drop() error {
	return os.Remove(f.path)
}

type mptIterator struct {
	kvdb.Iterator
}

func (it mptIterator) Next() bool {
	for it.Iterator.Next() {
		if evmstore.IsMptKey(it.Key()) {
			return true
		}
	}
	return false
}

type mptAndPreimageIterator struct {
	kvdb.Iterator
}

func (it mptAndPreimageIterator) Next() bool {
	for it.Iterator.Next() {
		if evmstore.IsMptKey(it.Key()) || evmstore.IsPreimageKey(it.Key()) {
			return true
		}
	}
	return false
}

type unitWriter struct {
	plain            io.WriteSeeker
	gziper           *gzip.Writer
	fileshasher      *fileshash.Writer
	dataStartPos     int64
	uncompressedSize uint64
}

func newUnitWriter(plain io.WriteSeeker) *unitWriter {
	return &unitWriter{
		plain: plain,
	}
}

func (w *unitWriter) Start(header genesis.Header, name, tmpDirPath string) error {
	if w.plain == nil {
		// dry run
		w.fileshasher = fileshash.WrapWriter(nil, genesisstore.FilesHashPieceSize, func(int) fileshash.TmpWriter {
			return devnullfile.DevNull{}
		})
		return nil
	}
	// Write unit marker and version
	_, err := w.plain.Write(append(genesisstore.FileHeader, genesisstore.FileVersion...))
	if err != nil {
		return err
	}

	// write genesis header
	err = rlp.Encode(w.plain, genesisstore.Unit{
		UnitName: name,
		Header:   header,
	})
	if err != nil {
		return err
	}

	w.dataStartPos, err = w.plain.Seek(8+8+32, io.SeekCurrent)
	if err != nil {
		return err
	}

	w.gziper, _ = gzip.NewWriterLevel(w.plain, gzip.BestCompression)

	w.fileshasher = fileshash.WrapWriter(w.gziper, genesisstore.FilesHashPieceSize, func(tmpI int) fileshash.TmpWriter {
		tmpI++
		tmpPath := path.Join(tmpDirPath, fmt.Sprintf("genesis-%s-tmp-%d", name, tmpI))
		_ = os.MkdirAll(tmpDirPath, os.ModePerm)
		tmpFh, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Crit("File opening error", "path", tmpPath, "err", err)
		}
		return dropableFile{
			ReadWriteSeeker: tmpFh,
			Closer:          tmpFh,
			path:            tmpPath,
		}
	})
	return nil
}

func (w *unitWriter) Flush() (hash.Hash, error) {
	if w.plain == nil {
		return w.fileshasher.Root(), nil
	}
	h, err := w.fileshasher.Flush()
	if err != nil {
		return hash.Hash{}, err
	}

	err = w.gziper.Close()
	if err != nil {
		return hash.Hash{}, err
	}

	endPos, err := w.plain.Seek(0, io.SeekCurrent)
	if err != nil {
		return hash.Hash{}, err
	}

	_, err = w.plain.Seek(w.dataStartPos-(8+8+32), io.SeekStart)
	if err != nil {
		return hash.Hash{}, err
	}

	_, err = w.plain.Write(h.Bytes())
	if err != nil {
		return hash.Hash{}, err
	}
	_, err = w.plain.Write(bigendian.Uint64ToBytes(uint64(endPos - w.dataStartPos)))
	if err != nil {
		return hash.Hash{}, err
	}
	_, err = w.plain.Write(bigendian.Uint64ToBytes(w.uncompressedSize))
	if err != nil {
		return hash.Hash{}, err
	}

	_, err = w.plain.Seek(0, io.SeekEnd)
	if err != nil {
		return hash.Hash{}, err
	}
	return h, nil
}

func (w *unitWriter) Write(b []byte) (n int, err error) {
	n, err = w.fileshasher.Write(b)
	w.uncompressedSize += uint64(n)
	return
}

func exportGenesis(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	fn := ctx.Args().First()

	from := idx.Epoch(1)
	if len(ctx.Args()) > 1 {
		n, err := strconv.ParseUint(ctx.Args().Get(1), 10, 32)
		if err != nil {
			return err
		}
		from = idx.Epoch(n)
	}
	to := idx.Epoch(math.MaxUint32)
	if len(ctx.Args()) > 2 {
		n, err := strconv.ParseUint(ctx.Args().Get(2), 10, 32)
		if err != nil {
			return err
		}
		to = idx.Epoch(n)
	}
	mode := ctx.String(EvmExportMode.Name)
	if mode != "full" && mode != "ext-mpt" && mode != "mpt" && mode != "none" {
		return errors.New("--export.evm.mode must be one of {full, ext-mpt, mpt, none}")
	}

	return ExportGenesis(ctx, fn, from, to, mode)
}

func ExportGenesis(ctx *cli.Context, fn string, from idx.Epoch, to idx.Epoch, mode string) error {
	cfg := makeAllConfigs(ctx)
	tmpPath := path.Join(cfg.Node.DataDir, "tmp")
	_ = os.RemoveAll(tmpPath)
	defer os.RemoveAll(tmpPath)

	rawDbs := makeDirectDBsProducer(cfg)
	gdb := makeGossipStore(rawDbs, cfg)
	if gdb.GetHighestLamport() != 0 {
		log.Warn("Attempting genesis export not in a beginning of an epoch. Genesis file output may contain excessive data.")
	}
	defer gdb.Close()

	// Open the file handle
	var plain io.WriteSeeker
	if fn != "dry-run" {
		fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer fh.Close()
		plain = fh
	}

	header := genesis.Header{
		GenesisID:   *gdb.GetGenesisID(),
		NetworkID:   gdb.GetEpochState().Rules.NetworkID,
		NetworkName: gdb.GetEpochState().Rules.Name,
	}
	var epochsHash hash.Hash
	var blocksHash hash.Hash
	var evmHash hash.Hash

	if from < 1 {
		// avoid underflow
		from = 1
	}
	if to > gdb.GetEpoch() {
		to = gdb.GetEpoch()
	}
	toBlock := idx.Block(0)
	fromBlock := idx.Block(0)
	{
		log.Info("Exporting epochs", "from", from, "to", to)
		writer := newUnitWriter(plain)
		err := writer.Start(header, genesisstore.EpochsSection, tmpPath)
		if err != nil {
			return err
		}
		for i := to; i >= from; i-- {
			er := gdb.GetFullEpochRecord(i)
			if er == nil {
				log.Warn("No epoch record", "epoch", i)
				break
			}
			b, _ := rlp.EncodeToBytes(ier.LlrIdxFullEpochRecord{
				LlrFullEpochRecord: *er,
				Idx:                i,
			})
			_, err := writer.Write(b)
			if err != nil {
				return err
			}
			if i == from {
				fromBlock = er.BlockState.LastBlock.Idx
			}
			if i == to {
				toBlock = er.BlockState.LastBlock.Idx
			}
		}
		epochsHash, err = writer.Flush()
		if err != nil {
			return err
		}
		log.Info("Exported epochs", "hash", epochsHash.String())
	}

	if fromBlock < 1 {
		// avoid underflow
		fromBlock = 1
	}
	{
		log.Info("Exporting blocks", "from", fromBlock, "to", toBlock)
		writer := newUnitWriter(plain)
		err := writer.Start(header, genesisstore.BlocksSection, tmpPath)
		if err != nil {
			return err
		}
		for i := toBlock; i >= fromBlock; i-- {
			br := gdb.GetFullBlockRecord(i)
			if br == nil {
				log.Warn("No block record", "block", i)
				break
			}
			if i%200000 == 0 {
				log.Info("Exporting blocks", "last", i)
			}
			b, _ := rlp.EncodeToBytes(ibr.LlrIdxFullBlockRecord{
				LlrFullBlockRecord: *br,
				Idx:                i,
			})
			_, err := writer.Write(b)
			if err != nil {
				return err
			}
		}
		blocksHash, err = writer.Flush()
		if err != nil {
			return err
		}
		log.Info("Exported blocks", "hash", blocksHash.String())
	}

	if mode != "none" {
		log.Info("Exporting EVM data", "from", fromBlock, "to", toBlock)
		writer := newUnitWriter(plain)
		err := writer.Start(header, genesisstore.EvmSection, tmpPath)
		if err != nil {
			return err
		}
		it := gdb.EvmStore().EvmDb.NewIterator(nil, nil)
		if mode == "mpt" {
			// iterate only over MPT data
			it = mptIterator{it}
		} else if mode == "ext-mpt" {
			// iterate only over MPT data and preimages
			it = mptAndPreimageIterator{it}
		}
		defer it.Release()
		err = iodb.Write(writer, it)
		if err != nil {
			return err
		}
		evmHash, err = writer.Flush()
		if err != nil {
			return err
		}
		log.Info("Exported EVM data", "hash", evmHash.String())
	}

	fmt.Printf("- Genesis ID: %v \n", header.GenesisID.String())
	fmt.Printf("- Epochs hash: %v \n", epochsHash.String())
	fmt.Printf("- Blocks hash: %v \n", blocksHash.String())
	fmt.Printf("- EVM hash: %v \n", evmHash.String())

	return nil
}
