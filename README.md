# Loom&Doom (dogmatix)

Projektarbete på kursen [Operativsystem och processorienterad programmering
(1DT096) våren 2020][homepage], [Uppsala universitet][uu].

[homepage]: https://www.it.uu.se/education/course/homepage/os/vt19/project/

[uu]: https://www.uu.se/

Projektet kan beskrivas.

## Vad som krävs

För att köra krävs att openGL och Pixel är installerat för Go.

## Pixel:  
```go get github.com/faiface/pixel```

```go get github.com/faiface/glhf```

## openGL: 

```go get -u github.com/go-gl/glfw/v3.3/glfw```

- På macOS behövs Xcode eller Command Line Tools för Xcode (xcode-select --install) för headers och libraries.
- På Ubuntu/Debian-liknande Linux distributioner behövs `libgl1-mesa-dev och xorg-dev` packages.
- På CentOs/Fedora-liknande distributioner behövs `libX11-devel libXcursor-devel libXrandr-devel linXinerama-devel mesa-libGL-devel libXi-devel` packages.

Titta [här](https://www.glfw.org/docs/latest/compile.html#compile_deps) om du använder annat OS.


## Kom igång

För att starta, starta projektet.

För att bygga, bygg projektet.

## Katalogstruktur

client
game
server
