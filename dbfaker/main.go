package main

import "C"
import (
    "os"
    "fmt"
    "path"
    "log"
    "context"
    // "math/big"

    "github.com/ledgerwatch/erigon/core/types/accounts"
    "github.com/ledgerwatch/erigon/common"
    // "github.com/ledgerwatch/erigon/core/rawdb"
    "github.com/ledgerwatch/erigon/core/state"
    "github.com/ledgerwatch/erigon-lib/kv/mdbx"
    "github.com/ledgerwatch/erigon-lib/kv"
	ledgerLog "github.com/ledgerwatch/log/v3"
)

func main() {}

// This dir should be cleaned between each run
const DB_DIR = "build/chaindata"

//export PutAccount
func PutAccount(address []byte, rlpAccount []byte, incarnation uint64) {
    db := dbOpen()
    defer dbClose(db)

    var acct accounts.Account
    if err := acct.DecodeForHashing(rlpAccount); err != nil {
        log.Fatal(err)
    }
    acct.Incarnation = incarnation

    ctx := context.Background()
    tx, err := db.BeginRw(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer tx.Rollback()

    w := state.NewPlainStateWriterNoHistory(tx)
    w.UpdateAccountData(common.BytesToAddress(address), new(accounts.Account), &acct)

    if err = tx.Commit(); err != nil {
        log.Fatal(err)
    }
}

func dbOpen() kv.RwDB {
    logger := ledgerLog.New()
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    db, err := mdbx.NewMDBX(logger).Path(path.Join(cwd, DB_DIR)).Open()
    if err != nil {
        log.Fatal(err)
    }

    return db
}

func dbClose(db kv.RwDB) {
    db.Close()
    fmt.Print("Go mdbx closed\n")
}
