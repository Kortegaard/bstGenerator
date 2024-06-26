package functionBuilder


func FormatAuthors(authorFormat string) string{
    var l = `FUNCTION {write.authors}{
    author 's :=
    s num.names$ 'numnames :=
    #1 'nameptr :=
    { nameptr numnames #1 + < }
    {
        s nameptr "`+ authorFormat +`" format.name$
        write$

        nameptr numnames #1 - =    % if last number
        {" and " write$}
        {", " write$}
        if$
        nameptr #1 + 'nameptr :=
    }
    while$
} `
    return l

}

func IsPreprint() string{
    return`FUNCTION {is.arxiv} {
    #0
    % publisher is arxiv
    publisher empty$ { #0 } { "arxiv" publisher "l" change.case$ = } if$ or
    % journal name starts with arxiv
    journal empty$ { #0 } { "arxiv" journal "l" change.case$ #1 #5 substring$ = } if$ or
}
FUNCTION {article}
{
    is.arxiv
    { article.arxiv }
    { article.published }
    if$
}`
}

func WriteBibitem() string{
    return `FUNCTION{write.bibitem}{
    "\bibitem[" write$
    label write$
    "]{" write$
    cite$ write$
    "}" write$
    newline$
}`
}

func BasicFunctions() string{
return `FUNCTION {not}
{   { #0 }
    { #1 }
  if$
}

FUNCTION {and}
{   'skip$
    { pop$ #0 }
  if$
}

FUNCTION {or}
{   { pop$ #1 }
    'skip$
  if$
}`
}



func WriteLabelConstructor() string{
    return `FUNCTION {construct.label}{
    author 's :=
    s num.names$ 'numnames :=

    numnames #1 >
    { %% Multiple names
        "" % empty string to concat with what is added on stack
        #1 'nameptr :=
        { nameptr numnames #1 + < }
        {
            s nameptr "{v{}}{l{}}" format.name$
            *
            nameptr #1 + 'nameptr :=
        }
        while$
        'label :=
    }{ %% One name 
        s #1 "{ll}" format.name$
        #1 #3 substring$
        'label :=
    }
    if$
    %% Could year 
    label year
    #3 #2 substring$
    * 'label :=
}`
}

func BibBegin() string{
    return `FUNCTION {bib.begin}
{
    "\begin{thebibliography}{" longest.label "}" * * write$ newline$
    "\newcommand*{\doi}[1]{\href{https://doi.org/#1}{\sloppy #1}}" write$ 
    newline$
    newline$
} `
}
func BibEnd() string{
    return `FUNCTION {bib.end}
{
    "\end{thebibliography}" write$ newline$
} `
}


func Sorter() string{
    return `FUNCTION {format.authors.sort}{
    author 's :=
    s num.names$ 'numnames :=
    ""
    #1 'nameptr :=
    { nameptr numnames #1 + < }
    {
        s nameptr "{ll}{ff}" format.name$
        *
        nameptr numnames #1 - =    % if last number
        {" and " *}
        {", " *}
        if$
        nameptr #1 + 'nameptr :=
    }
    while$
} 

FUNCTION {presort}
{
    format.authors.sort 
    "     "  *
    year *
    "     "  *
    title *
    'sort.key$ :=
}

ITERATE {presort}`
}


func LabelDifferentiator() string{
    return `
    STRINGS {last.label} 
INTEGERS {addon.chr.track last.label.had.addon}

FUNCTION {init.label.passes}{
    "" 'last.label :=
    #0 'addon.chr.track :=
    #0 'last.label.had.addon :=
}

FUNCTION {label.forward} {
    label last.label =
    {
        label 'last.label :=
        label addon.chr.track "b" chr.to.int$ + int.to.chr$  * 'label :=
        addon.chr.track #1 + 'addon.chr.track  :=
        #1 'label.addon :=

    }{
        label 'last.label :=
        #0 'addon.chr.track :=
        #0 'label.addon :=
    }
    if$
}

FUNCTION {label.backwards} {
    #1 last.label.had.addon 
    label.addon #1 < and =
    {
        label "a" * 'label :=
        #1 'label.addon :=
    }{ }
    if$
    label.addon 'last.label.had.addon :=
}

EXECUTE {init.label.passes}
ITERATE {label.forward}
REVERSE {label.backwards}`
}

func FindLongestLabel() string{
    return `INTEGERS {longest.label.width}
STRINGS {longest.label}

FUNCTION {initiate.longest.label}{
    #0 'longest.label.width :=
    "" 'longest.label :=
}

FUNCTION {longest.label.search} {
    label text.length$ longest.label.width >
    {
        label 'longest.label :=
        label text.length$ 'longest.label.width :=
    }{}
    if$
}

EXECUTE {initiate.longest.label}
ITERATE {longest.label.search}`

}
