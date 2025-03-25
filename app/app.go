// package app
// func initApp(){
// 	ChainClient, err = client.NewClient(
// 		&c_config.ClientConfig{
// 			Version_:                cconfig.MetaNodeVersion,
// 			PrivateKey_:             cconfig.PrivateKey_,
// 			ParentAddress:           cconfig.ParentAddress,
// 			ParentConnectionAddress: cconfig.NodeConnectionAddress,
// 			DnsLink_:                cconfig.DnsLink(),
// 		},
// 	)

// 	log.Println("Config ok")

// 	// create card abi
// 	reader, err := os.Open(cconfig.MailFactoryABIPath)
// 	if err != nil {
// 		log.Fatalf("Error occured while read baccarat abi")
// 	}
// 	defer reader.Close()

// 	mailFactoryAbi, err := abi.JSON(reader)
// 	if err != nil {
// 		log.Fatalf("Error occured while parse baccarat smart contract abi")
// 	}
// 	//
// 	readerMailStorage, err := os.Open(cconfig.MailStorageABIPath)
// 	if err != nil {
// 		log.Fatalf("Error occured while read baccarat abi")
// 	}
// 	defer readerMailStorage.Close()

// 	abiMailStorage, err := abi.JSON(readerMailStorage)
// 	if err != nil {
// 		log.Fatalf("Error occured while parse baccarat smart contract abi")
// 	}
// 	//
// 	readerFileStorage, err := os.Open(cconfig.FileABIPath)
// 	if err != nil {
// 		log.Fatalf("Error occured while read baccarat abi")
// 	}
// 	defer readerMailStorage.Close()

// 	abiFile, err := abi.JSON(readerFileStorage)
// 	if err != nil {
// 		log.Fatalf("Error occured while parse baccarat smart contract abi")
// 	}
// 	//
// 	servs := services.NewSendTransactionService(
// 		ChainClient,
// 		&mailFactoryAbi,
// 		common.HexToAddress(cconfig.MailFactoryAddress),
// 		&abiMailStorage,
// 		common.HexToAddress(cconfig.MailStorageAddress),
// 		common.HexToAddress(cconfig.NotiAddress),
// 		common.HexToAddress(cconfig.AdminAddress),
// 		&abiFile,
// 		common.HexToAddress(cconfig.FileAddress),
// 	)

// }