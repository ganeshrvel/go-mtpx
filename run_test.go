package mtpx

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func testRun(t *testing.T) {
	//dev, err := Initialize(Init{DebugMode: true})
	dev, err := Initialize(Init{DebugMode: false})
	if err != nil {
		log.Panic(err)
	}

	_, err = FetchDeviceInfo(dev)
	if err != nil {
		log.Panic(err)
	}

	//pretty.Println(inf)

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Printf("outer >> error: %+v:", err)
		//Dispose(dev)

		time.Sleep(3 * time.Second)

		dev, err := Initialize(Init{DebugMode: true})
		if err != nil {
			log.Panic(err)
		}

		storages, err := FetchStorages(dev)
		if err != nil {
			log.Printf("inner >> error: %+v:", err)

			return
		}

		sid := storages[0].Sid
		log.Printf("inner >> storage id: %+v\n", sid)

		return
	}

	sid := storages[0].Sid
	//log.Printf("outer >> storage id: %+v\n", sid)
	//////////////////////
	//////////////////////
	//////////////////////
	//copy mock test files

	//uploadFile1 := getTestMocksAsset("a.txt")
	uploadFile2 := getTestMocksAsset("")
	//uploadFile3 := getTestMocksAsset("4mb_txt_file_2")
	sources := []string{uploadFile2}
	destination := "/test-mtp"
	_, _, _, err = UploadFiles(dev, sid,
		sources,
		destination,
		true,
		//func(fi *os.FileInfo, err error) error {
		//	if err != nil {
		//		return err
		//	}
		//
		//	fmt.Printf("Preprocessing File name: %s\n", (*fi).Name())
		//
		//	return nil
		//},
		func(fi *os.FileInfo, fullPath string, err error) error {
			if err != nil {
				return err
			}

			fmt.Printf("Preprocessing File name: %s\n", (*fi).Name())

			return nil
		},
		func(pi *ProgressInfo, err error) error {
			fmt.Printf("\nFile name: %s\n", pi.FileInfo.FullPath)
			//fmt.Printf("Total size: %d\n", pi.ActiveFileSize.Total)
			fmt.Printf("Size sent: %d\n", pi.ActiveFileSize.Sent)
			//fmt.Printf("Speed: %f\n", pi.Speed)
			//fmt.Printf("Object Id: %d\n", pi.FileInfo.ObjectId)
			//fmt.Printf("ActiveFileSize progress: %f\n", pi.ActiveFileSize.Progress)
			//fmt.Printf("TotalFiles: %d\n", pi.TotalFiles)
			//fmt.Printf("totalDirectories: %d\n", pi.TotalDirectories)
			//fmt.Printf("FilesSent: %d\n", pi.FilesSent)
			//fmt.Printf("FilesSentProgress: %f\n\n\n", pi.FilesSentProgress)

			return nil
		},
	)
	if err != nil {
		log.Panicln(err)
	}
	//filePath := getTestMocksAsset("mock_dir1/[-----DS_Store.mtp.test----].txt")
	//Exists, _ := FileExists(dev, sid, 0, filePath)
	//if !Exists {
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
	//	FileInfo, err := os.Lstat(filePath)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	fileBuf, err := os.Open(filePath)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	_, err = handleMakeFile(dev, sid, &fObj, &FileInfo, fileBuf, true, func(total, sent int64, objectId uint32, err error) error {
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
	//Exists := FileExists(dev, Sid, 0, "/tests/test.txt")
	//
	//pretty.Println("======\n")
	//pretty.Println("Does File Exists:", Exists)

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
	//sources := []string{uploadFile1}
	//destination := "/mtp-test-files/mock_dir1"
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
	//		fmt.Printf("Total size: %d\n", pi.ActiveFileSize.Total)
	//		fmt.Printf("Size sent: %d\n", pi.ActiveFileSize.Sent)
	//		fmt.Printf("Speed: %f\n", pi.Speed)
	//		//fmt.Printf("Object Id: %d\n", pi.FileInfo.ObjectId)
	//		fmt.Printf("ActiveFileSize progress: %f\n", pi.ActiveFileSize.Progress)
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

	//sourceFile1 := "/mtp-test-files/mock_dir1/a.txt"
	//sourceFile2 := "/mtp-test-files/mock_dir1/"
	//destination := newTempMocksDir("test_DownloadTest", true)
	//
	//_, _, err = DownloadFiles(dev, sid,
	//	[]string{sourceFile1, sourceFile2}, destination, true,
	//	func(fi *FileInfo, err error) error {
	//		fmt.Printf("Preprocessing files 'FullPath': %s\n", fi.FullPath)
	//		fmt.Printf("Preprocessing files 'Size': %d\n", fi.Size)
	//
	//		return nil
	//	},
	//	func(fi *ProgressInfo, err error) error {
	//		fmt.Printf("Current filepath: %s\n", fi.FileInfo.FullPath)
	//		fmt.Printf("%f MB/s\n", fi.Speed)
	//		fmt.Printf("BulkFileSize Total: %d\n", fi.BulkFileSize.Total)
	//		fmt.Printf("BulkFileSize Sent: %d\n", fi.BulkFileSize.Sent)
	//		fmt.Printf("ActiveFileSize Total: %d\n", fi.ActiveFileSize.Total)
	//		fmt.Printf("ActiveFileSize Sent: %d\n\n", fi.ActiveFileSize.Sent)
	//
	//		return nil
	//	},
	//)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//Dispose(dev)
}
