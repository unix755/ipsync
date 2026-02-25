package preload

import (
	"fmt"
	"os"
)

func (p Preload) SaveToFileOrPrint(savedFilepath string, encryptionKey []byte) error {
	bytes, err := Marshal(p, "json", encryptionKey)
	if err != nil {
		return err
	}

	// 保存到文件或者打印
	if savedFilepath != "" {
		return os.WriteFile(savedFilepath, bytes, 0644)
	}
	fmt.Println(string(bytes))
	return nil
}
