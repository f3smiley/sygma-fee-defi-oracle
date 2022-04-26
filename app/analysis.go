// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package app

import "net/http"
import _ "net/http/pprof" // #nosec

func (a *FeeOracleApp) goroutineMemoryLeakChecker() {
	if a.base.GetEnv() != AppEvmProd {
		go func() {
			a.base.GetLogger().Info(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}
