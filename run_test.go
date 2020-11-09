package mtpx

import (
	"github.com/kr/pretty"
	"log"
	"testing"
)

func TestRun(t *testing.T) {
	dev, err := Initialize(Init{DebugMode: false})

	if err != nil {
		log.Panic(err)
	}

	_, err = FetchDeviceInfo(dev)
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].Sid
	pretty.Println("storage id: ", sid)

	////////////////////////
	////////////////////////
	////////////////////////
	// copy mock test files
	//uploadFile1 := getTestMocksAsset("")
	//sources := []string{uploadFile1}
	//destination := "/"
	//_, _, _, err = UploadFiles(dev, sid,
	//	sources,
	//	destination,
	//	true,
	//	func(fi *os.FileInfo, err error) error {
	//		if err != nil {
	//			return err
	//		}
	//
	//		fmt.Printf("Preprocessing File name: %s\n", (*fi).Name())
	//
	//		return nil
	//	},
	//	func(pi *ProgressInfo, err error) error {
	//		//fmt.Printf("\nFile name: %s\n", pi.FileInfo.FullPath)
	//		//fmt.Printf("Total size: %d\n", pi.ActiveFileSize.Total)
	//		//fmt.Printf("Size sent: %d\n", pi.ActiveFileSize.Sent)
	//		//fmt.Printf("Speed: %f\n", pi.Speed)
	//		//fmt.Printf("Object Id: %d\n", pi.FileInfo.ObjectId)
	//		//fmt.Printf("ActiveFileSize progress: %f\n", pi.ActiveFileSize.Progress)
	//		//fmt.Printf("TotalFiles: %d\n", pi.TotalFiles)
	//		//fmt.Printf("totalDirectories: %d\n", pi.TotalDirectories)
	//		//fmt.Printf("FilesSent: %d\n", pi.FilesSent)
	//		//fmt.Printf("FilesSentProgress: %f\n\n\n", pi.FilesSentProgress)
	//
	//		return nil
	//	},
	//)
	//if err != nil {
	//	log.Panicln(err)
	//}
	//filePath := getTestMocksAsset("mock_dir1/[-----DS_Store.mtp.test----].txt")
	//exists, _ := FileExists(dev, sid, 0, filePath)
	//if !exists {
	//	destination := "/mtp-test-files/mock_dir1"
	//	fName := "[-----DS_Store.mtp.test----].txt"
	//
	//	dFi, err := GetObjectFromObjectIdOrPath(dev, sid, 0, destination)
	//
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//
	//	fObj := mtp.ObjectInfo{
	//		StorageID:        sid,
	//		ObjectFormat:     mtp.OFC_Undefined,
	//		ParentObject:     dFi.ObjectId,
	//		Filename:         fName,
	//		CompressedSize:   0,
	//		ModificationDate: time.Now(),
	//	}
	//
	//	fileInfo, err := os.Lstat(filePath)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	fileBuf, err := os.Open(filePath)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_, err = handleMakeFile(dev, sid, &fObj, &fileInfo, fileBuf, true, func(total, sent int64, objectId uint32, err error) error {
	//
	//		return nil
	//	})
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

	////////////////////////
	////////////////////////
	////////////////////////

	//totalFiles, err := dev.GetNumObjects(Sid, mtp.GOH_ALL_ASSOCS, ParentObjectId)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(int64(totalFiles))
	//

	//objectId, totalFiles, err := Walk(dev, sid, "/mtp-test-files/mock_dir1", true, false, func(objectId uint32, fi *FileInfo, err error) error {
	//	pretty.Println("Filepath ", (*fi).FullPath)
	//
	//	return nil
	//})
	//
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println("totalFiles: ", totalFiles)
	//pretty.Println("objectId: ", objectId)

	////MakeDirectory
	//objectId, err := MakeDirectory(dev, Sid, ParentObjectId, "/", "name")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(objectId)

	//GetObjectFromPath
	//fileObj, err := GetObjectFromPath(dev, Sid, "/tests/s")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println("======\n")
	//pretty.Println(fileObj)
	//

	// FileExists
	//exists := FileExists(dev, Sid, 0, "/tests/test.txt")
	//
	//pretty.Println("======\n")
	//pretty.Println("Does File exists:", exists)

	///DeleteFile
	//err = DeleteFile(dev, Sid, 0, "/mtp-test-files/temp_dir/this is a test")
	//if err != nil {
	//	log.Panic(err)
	//}

	//////RenameFile
	//objId, err := RenameFile(dev, sid, 0, "/mtp-test-files/mock_dir1/.DS_Store.txt", ".DS_Store")
	//if err != nil {
	//	log.Panic(err)
	//}
	//pretty.Println(objId)

	//UploadFiles
	//start := time.Now()
	//uploadFile1 := getTestMocksAsset("")
	////uploadFile2 := getTestMocksAsset("mock_dir2")
	//sources := []string{uploadFile1}
	//destination := "/mtp-test-files/temp_dir/test_UploadFiles"
	//_, _, _, err = UploadFiles(dev, sid,
	//	sources,
	//	destination,
	//	true,
	//	func(fi *os.FileInfo, err error) error {
	//		if err != nil {
	//			return err
	//		}
	//
	//		fmt.Printf("Preprocessing File name: %s\n", (*fi).Name())
	//
	//		return nil
	//	},
	//	func(pi *ProgressInfo, err error) error {
	//		fmt.Printf("File name: %s\n", pi.FileInfo.FullPath)
	//		//fmt.Printf("Total size: %d\n", pi.ActiveFileSize.Total)
	//		//fmt.Printf("Size sent: %d\n", pi.ActiveFileSize.Sent)
	//		//fmt.Printf("Speed: %f\n", pi.Speed)
	//		//fmt.Printf("Object Id: %d\n", pi.FileInfo.ObjectId)
	//		//fmt.Printf("ActiveFileSize progress: %f\n", pi.ActiveFileSize.Progress)
	//		fmt.Printf("TotalFiles: %d\n", pi.TotalFiles)
	//		fmt.Printf("totalDirectories: %d\n", pi.TotalDirectories)
	//		fmt.Printf("FilesSent: %d\n", pi.FilesSent)
	//		fmt.Printf("FilesSentProgress: %f\n\n\n", pi.FilesSentProgress)
	//
	//		return nil
	//	},
	//)

	//pretty.Println(objectIdDest)
	//pretty.Println(totalFiles)
	//pretty.Println(totalSize)
	//pretty.Println("time elapsed: ", time.Since(start).Seconds())

	//
	//totalFiles, totalSize, err := DownloadFiles(dev, Sid,
	//	[]string{sourceFile1}, downloadFile,
	//	func(downloadFi *TransferredFileInfo, err error) error {
	//		fmt.Printf("ActiveFileSize filepath: %s\n", downloadFi.FileInfo.FullPath)
	//		fmt.Printf("%f MB/s\n", downloadFi.Speed)
	//
	//		return nil
	//	},
	//)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(totalFiles)
	//pretty.Println(totalSize)
	//pretty.Println("time elapsed: ", time.Since(start).Seconds())

	Dispose(dev)
}
