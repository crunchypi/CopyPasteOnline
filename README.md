# CopyPasteOnline
Source code for a micro service for doing online copy-paste using mnemonics

## API usage:
Start server in root/backend/src.main.go
This file contains a simple upspin:

    // # Setup db manager
	  db, err := sqlite.New("./data.sqlite")
	  if err != nil {
		  panic("unable to setup db: " + err.Error())
	  }

	  // # Setup mnemonics handler Will panic on error.
	  mnemonics := mnemonics.New()

	  // # Stat server.
	  app.Start(db, mnemonics)
    
If started with default configs, requests can be made as such (pseudocode):
    
    // # Store data, this will return a mnemonic string on success.
    POST https://localhost:8080:/transfer -data "demonstration data"
    
    // # Get data, this will return the stored data, given a mnemonic string.
    GET https://localhost:8080:/transfer -mnemonic "ice tea pie"
    

Curl snippets:

    // # Store data, this will return a mnemonic string on success. (-k allows unsafe cert)
    curl --request POST https://https://localhost:8080:/transfer -d "demonstration data" -k
    
    // # Get data, this will return the stored data, given a mnemonic string.
    curl --request GET https://localhost:8080/transfer -d "ice tea pie" -k -tls
    
