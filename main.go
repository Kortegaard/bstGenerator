package main

import(
    "os"
    "fmt" 
    "strings"
    "bst_generator/functionBuilder"
) 


type EntryTokenCollection []Token

type TokenType int64

type BstReturnType struct{
    inScope string
    outOfScope string
}

type Token interface{
    constructBst() BstReturnType
}

type TokenText struct{
    text string
}

func (t TokenText) constructBst() BstReturnType{
    var res string = "\""
    res += t.text
    res += "\""
    res += " write$"
    return BstReturnType{ inScope: res, outOfScope: "" }
}

type TokenVariable struct{
    optional bool
    format string
    variableName string
    preText string
    postText string
}

func (t TokenVariable) constructBst() BstReturnType{
    var inScope string = ""
    if t.preText != ""{
        inScope += "\"" + t.preText + "\" write$\n"
    }
    switch t.variableName{
        case "author":
            inScope += "write.authors"
        default:
            inScope += t.variableName
            inScope += " write$"
    }
    if t.postText != ""{
        inScope += "\n\"" + t.postText + "\" write$"
    }
    
    if t.optional{
        res := "\n"
        res += t.variableName + " empty$ {}\n    {\n"
        res += "        " + inScope 
        res += "\n    }\n    if$\n"
        inScope = res
    } 

    return BstReturnType{ inScope: inScope, outOfScope: "" }
}




func main() {
    b := initBaseBstBuilder()
    b.addEntryFromFormat("article.arxiv", "[[author:{f}{ll}]], {\\it [[title]]}, preprint ([[year]]). doi:[[doi]].")
    b.addEntryFromFormat("article.published", "[[author:{f}{ll}]], {\\it [[title]]}, [[journal]] {\\bf [[volume]]} ([[year]])[, no. [?number],] [[pages]]. [ doi:\\doi{[?doi]}]")
    b.addEntryFromFormat("book", "[[author:{f}{ll}]], {\\it [[title]]}, [[journal]] {\\bf [[volume]]} ([[year]])[, no. [?number],] [[pages]].d[ doi:\\doi{[?doi]}.]")
    st := b.build()
    f, err := os.Create("./mine.bst")
    if err != nil{
        panic(err)
    }
    defer f.Close()
    f.WriteString(st)

    //fmt.Println(b.build())
    //fmt.Println(functionBuilder.FormatAuthors())
}

func ParseVariableEnvironment(s string) Token{
    var token TokenVariable

    var i int = 0
    var j int = 0
    i,j = FindNextBracketPair(s)
    if i == -1{ return token } // TODO: add this as an error instead
    token.preText = s[0:i]
    token.postText = s[j+1:]
    ParseVariable(s[i+1:j], &token)
    return token
}

func ParseVariable(s string, token *TokenVariable){
    var variableNameStartIndex int = 0 
    var variableNameEndIndex int = len(s)
    token.optional = (s[0] == '?')
    if token.optional{
        variableNameStartIndex += 1
    }

    var colonIndex int = strings.Index(s,":")
    if colonIndex != -1{
        variableNameEndIndex = colonIndex
        token.format = s[colonIndex+1:]
    }
    token.variableName = s[variableNameStartIndex:variableNameEndIndex]

    if colonIndex != -1 && token.variableName != "author"{
        fmt.Printf("WARNING: YOU FUCKED UP")
    }
}

func FindNextBracketPair(s string) (int, int){
    var sqDepth int = 0
    var sqStart int = 0
    for index, ch := range s{
        if ch == '['{
            if sqDepth == 0 {
                sqStart = index
            }
            sqDepth += 1
        }else if ch == ']'{
            sqDepth -= 1
            if sqDepth == 0 {
                return sqStart, index
            }
        }
    }
    return -1, -1
}

func ParseEntryFormat(s string)[]Token{
    var currIndex int = 0
    var i int = 0
    var j int = 0

    var tokens []Token

    for currIndex < len(s){
        i,j = FindNextBracketPair(s[currIndex:])
        if i == -1{
            tokens = append(tokens, TokenText{text: s[currIndex:]})
            break
        }
        i += currIndex
        j += currIndex
        tokens = append(tokens, TokenText{text: s[currIndex:i]})
        tokens = append(tokens, ParseVariableEnvironment(s[i+1:j]))
        currIndex = j + 1
    }
    return tokens
}

func initBaseBstBuilder() BSTBuilder{
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

func (b *BSTBuilder) addEntryFromFormat(entryName string, format string){
    var tokens []Token = ParseEntryFormat(format)

    var entryFunction string = "FUNCTION{" + entryName + "}{\n"
    entryFunction +="write.bibitem"
    for _, tok := range tokens{
        var ret BstReturnType = tok.constructBst()
        entryFunction += "    " + ret.inScope + "\n"
        if ret.outOfScope != ""{
            b.addCodeBeforeRead(ret.outOfScope)
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

func (b BSTBuilder) build() string{

    res := ""
    res += "ENTRY{\n"
    res += b.mBuildFields() + "\n"
    res += "}{}{label}\n\n"
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
    res += "SORT\n\n"
    res += b.mBuildAfterSort() + "\n\n"
    res += functionBuilder.BibBegin() + "\n\n"
    res += functionBuilder.BibEnd() + "\n\n"
    res += "EXECUTE {bib.begin}"
    res += "ITERATE {call.type$}" + "\n\n"
    res += "EXECUTE {bib.end}"

    return res
}



