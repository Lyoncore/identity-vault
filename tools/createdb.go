// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"log"

	"github.com/ubuntu-core/identity-vault/service"
)

func main() {
	env := service.Env{}
	// Parse the command line arguments
	service.ParseArgs()
	service.ReadConfig(&env.Config)

	// Open the connection to the local database
	env.DB = service.OpenSysDatabase(env.Config.Driver, env.Config.DataSource)

	// Create the keypair table, if it does not exist
	err := env.DB.CreateKeypairTable()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created the 'keypair' table.")

		// Create the test key (if the filesystem store is used)
		if env.Config.KeyStoreType == "filesystem" {
			// Create the test key as it is in the default filesystem keystore
			env.DB.PutKeypair(service.Keypair{AuthorityID: "System", KeyID: "61abf588e52be7a3"})
		}
	}

	// Create the model table, if it does not exist
	err = env.DB.CreateModelTable()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created the 'model' table.")
	}

	// Create the keypair table, if it does not exist
	err = env.DB.CreateSettingsTable()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Created the 'settings' table.")
	}

	// Initalize the TPM store, authenticating with the TPM 2.0 module
	if env.Config.KeyStoreType == service.TPM20Store.Name {
		log.Println("Initialize the TPM2.0 store")
		err = service.TPM2InitializeKeystore(env, nil)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Initialized TPM 2.0 module.")
		}
	}
}
