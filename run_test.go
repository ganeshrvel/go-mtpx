package mtpx

import (
	"fmt"
	"github.com/kr/pretty"
	"log"
	"os"
	"testing"
)

func _TestRun(t *testing.T) {
	//todo remove these
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

	//totalFiles, err := dev.GetNumObjects(Sid, mtp.GOH_ALL_ASSOCS, ParentObjectId)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(int64(totalFiles))
	//

	/*objectId, totalFiles, err := Walk(dev, Sid, 0, "/mtp-test-files/mock_dir1", true, func(objectId uint32, fi *FileInfo) {
		pretty.Println("objectId is: ", objectId)
	})

	if err != nil {
		log.Panic(err)
	}

	pretty.Println("totalFiles: ", totalFiles)
	pretty.Println("objectId: ", objectId)*/

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

	////RenameFile
	//objId, err := RenameFile(dev, Sid, 0, "/mtp-test-files/temp_dir/b.txt", "b.txt")
	//if err != nil {
	//	log.Panic(err)
	//}
	//pretty.Println(objId)

	//UploadFiles
	//start := time.Now()
	uploadFile1 := getTestMocksAsset("test-large-file")
	//uploadFile2 := getTestMocksAsset("mock_dir2")
	sources := []string{uploadFile1}
	destination := "/mtp-test-files/temp_dir/test_UploadFiles"
	_, _, _, err = UploadFiles(dev, sid,
		sources,
		destination,
		true,
		func(fi *os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fmt.Printf("Preprocessing File name: %s\n", (*fi).Name())

			return nil
		},
		func(pi *ProgressInfo, err error) error {
			fmt.Printf("File name: %s\n", pi.FileInfo.FullPath)
			//fmt.Printf("Total size: %d\n", pi.Current.Total)
			//fmt.Printf("Size sent: %d\n", pi.Current.Sent)
			//fmt.Printf("Speed: %f\n", pi.Speed)
			//fmt.Printf("Object Id: %d\n", pi.FileInfo.ObjectId)
			//fmt.Printf("Current progress: %f\n", pi.Current.Progress)
			fmt.Printf("TotalFiles: %d\n", pi.TotalFiles)
			fmt.Printf("totalDirectories: %d\n", pi.TotalDirectories)
			fmt.Printf("FilesSent: %d\n", pi.FilesSent)
			fmt.Printf("FilesSentProgress: %f\n\n\n", pi.FilesSentProgress)

			return nil
		},
	)

	//pretty.Println(objectIdDest)
	//pretty.Println(totalFiles)
	//pretty.Println(totalSize)
	//pretty.Println("time elapsed: ", time.Since(start).Seconds())

	//
	//totalFiles, totalSize, err := DownloadFiles(dev, Sid,
	//	[]string{sourceFile1}, downloadFile,
	//	func(downloadFi *TransferredFileInfo, err error) error {
	//		fmt.Printf("Current filepath: %s\n", downloadFi.FileInfo.FullPath)
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
