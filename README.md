# copypaste-api
Simple API for using copy-paste on the network, using mnemonics.




## API usage:
(configs at root/config/config.go)
<br>
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
    POST http://localhost:80:/copy -data "demonstration data"
    
    // # Get data, this will return the stored data, given a mnemonic string.
    POST http://localhost:80:/paste -mnemonic "ice tea pie"
    

Curl snippets:

    // # Store data, this will return a mnemonic string on success.
    curl --request POST http://localhost:80:/copy -d "demonstration data"
    
    // # Get data, this will return the stored data, given a mnemonic string.
    curl --request POST http://localhost:80/paste -d "ice tea pie"
    
