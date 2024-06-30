package main

import(
    "os"
    "github.com/Kortegaard/bstGenerator/bstBuilder"
) 


func main() {
    b := bstBuilder.InitBaseBstBuilder()
    b.AddEntryFromFormat("article.arxiv", "[[author:{f}{ll}]], {\\it [[title]]}, preprint ([[year]]). doi:[[doi]].")
    b.AddEntryFromFormat("article.published", "[[author:{f}{ll}]], {\\it [[title]]}, [[journal]] {\\bf [[volume]]} ([[year]])[, no. [?number],] [[pages]]. [ doi:\\doi{[?doi]}]")
    b.AddEntryFromFormat("book", "[[author:{f}{ll}]], {\\it [[title]]}, [[journal]] {\\bf [[volume]]} ([[year]])[, no. [?number],] [[pages]].d[ doi:\\doi{[?doi]}.]")
    st := b.Build()
    f, err := os.Create("./mine.bst")
    if err != nil{
        panic(err)
    }
    defer f.Close()
    f.WriteString(st)

    //fmt.Println(b.build())
    //fmt.Println(functionBuilder.FormatAuthors())
}

