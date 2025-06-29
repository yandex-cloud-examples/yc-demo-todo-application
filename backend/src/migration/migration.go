package migration
import (
    "os"
    "log"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "todo/model"
    "gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
    if err := db.AutoMigrate(&model.Todo{}); err != nil {
        log.Fatal(fmt.Sprintf("Migration failed: %s\n", err))
    }
    log.Print("Migration finished successfuly.\n")
}

func LoadData(db *gorm.DB, filename string) {
    var todos []model.Todo
    var first model.Todo

    if r := db.First(&first); r.Error == nil {
        log.Print("Target table isn't empty. Skip loading\n")
        return 
    }
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(fmt.Sprintf("Unable to open %s: %s\n", filename, err))
    }
    defer file.Close()
    buff, _ := ioutil.ReadAll(file)
    err = json.Unmarshal(buff, &todos)
    for _, todo := range todos {
        if r := db.Create(&todo); r.Error != nil {
            log.Fatal(fmt.Sprintf("Load form %s failed\n", filename))
        }
    }
    log.Print("Loading data finished successfuly.\n")
}
