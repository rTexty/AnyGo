package main

type Prototype interface{
    clone() Prototype
}

type File struct{
    name string
}

func (f *File) clone() Prototype{
    return &File{
        name: f.name,
    }
}

type Folder struct{
    ch []Prototype
    name string
}

func (f *Folder) clone() Prototype{
    var tempCh []Prototype
    for _, children := range f.ch{
        copy := children.clone()
        tempCh = append(tempCh, copy)
    }
    return &Folder{
        ch: tempCh,
        name: f.name,
    }
}
