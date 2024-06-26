package entryFormatter
import(
    "strings"
    "fmt"
) 

type EntryTokenCollection []Token

type TokenType int64

type BstReturnType struct{
    InScope string
    OutOfScope string
}

type Token interface{
    ConstructBst() BstReturnType
}

type TokenText struct{
    text string
}

func (t TokenText) ConstructBst() BstReturnType{
    var res string = "\""
    res += t.text
    res += "\""
    res += " write$"
    return BstReturnType{ InScope: res, OutOfScope: "" }
}

type TokenVariable struct{
    optional bool
    format string
    variableName string
    preText string
    postText string
}

func (t TokenVariable) ConstructBst() BstReturnType{
    var InScope string = ""
    if t.preText != ""{
        InScope += "\"" + t.preText + "\" write$\n"
    }
    switch t.variableName{
        case "author":
            InScope += "write.authors"
        default:
            InScope += t.variableName
            InScope += " write$"
    }
    if t.postText != ""{
        InScope += "\n\"" + t.postText + "\" write$"
    }
    
    if t.optional{
        res := "\n"
        res += t.variableName + " empty$ {}\n    {\n"
        res += "        " + InScope 
        res += "\n    }\n    if$\n"
        InScope = res
    } 

    return BstReturnType{ InScope: InScope, OutOfScope: "" }
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


