package bstBuilder

import(
    "strings"
    "bst_generator/functionBuilder"
    "bst_generator/entryFormatter"
) 

func InitBaseBstBuilder() BSTBuilder{
    var bstBuilder BSTBuilder
    bstBuilder.fields = []string{ "address", "author", "booktitle", "chapter", "doi", "edition", "editor", "howpublished", "institution", "journal", "key", "note", "number", "organization", "pages", "publisher", "school", "series", "title", "type", "volume", "year", }
    return bstBuilder
}

type BSTBuilder struct{
    authorFormat string
    fields []string
    beforeRead []string
    afterRead []string
    afterSort []string
}


func (b *BSTBuilder) addCodeBeforeRead(code string){
    b.beforeRead = append(b.beforeRead, code)
}

func (b *BSTBuilder) addCodeAfterRead(code string){
    b.afterRead = append(b.afterRead, code)
}

func (b *BSTBuilder) addCodeAfterSort(code string){
    b.afterSort = append(b.afterSort, code)
}

func (b *BSTBuilder) AddEntryFromFormat(entryName string, format string){
    var tokens []entryFormatter.Token = entryFormatter.ParseEntryFormat(format)

    var entryFunction string = "FUNCTION{" + entryName + "}{\n"
    entryFunction +="write.bibitem"
    for _, tok := range tokens{
        var ret entryFormatter.BstReturnType = tok.ConstructBst()
        entryFunction += "    " + ret.InScope + "\n"
        if ret.OutOfScope != ""{
            b.addCodeBeforeRead(ret.OutOfScope)
        }
    }
    entryFunction += "    newline$\n"
    entryFunction += "}"
    b.beforeRead = append(b.beforeRead , entryFunction)
}

func (b BSTBuilder) mBuildFields() string{
    return "    " + strings.Join(b.fields, "\n    ")
}

func (b BSTBuilder) mBuildBeforeRead() string{
    return strings.Join(b.beforeRead, "\n\n")
}

func (b BSTBuilder) mBuildAfterRead() string{
    return strings.Join(b.afterRead, "\n\n")
}

func (b BSTBuilder) mBuildAfterSort() string{
    return strings.Join(b.afterSort, "\n\n")
}

func (b BSTBuilder) Build() string{

    res := ""
    res += "ENTRY{\n"
    res += b.mBuildFields() + "\n"
    res += "}{label.addon}{label}\n\n"
    res += "STRINGS {s}\n\n"
    res += "INTEGERS {nameptr numnames}\n\n"

    res += functionBuilder.BasicFunctions() + "\n\n"
    res += functionBuilder.WriteLabelConstructor() + "\n\n"
    res += functionBuilder.FormatAuthors("{f. }{ll}") + "\n\n"
    res += functionBuilder.WriteBibitem() + "\n\n"
        
    res += b.mBuildBeforeRead() + "\n\n"

    res += functionBuilder.IsPreprint() + "\n\n"
    res += "READ\n\n"
    res += "ITERATE {construct.label}\n\n"
    res += b.mBuildAfterRead() + "\n\n"
    res += functionBuilder.Sorter() + "\n\n"
    res += "SORT\n\n"
    res += functionBuilder.LabelDifferentiator() + "\n\n"
    res += b.mBuildAfterSort() + "\n\n"
    res += functionBuilder.FindLongestLabel() + "\n\n"
    res += functionBuilder.BibBegin() + "\n\n"
    res += functionBuilder.BibEnd() + "\n\n"
    res += "EXECUTE {bib.begin}"
    res += "ITERATE {call.type$}" + "\n\n"
    res += "EXECUTE {bib.end}"

    return res
}



