## Contents

We list a few examples of the `Image Generator`'s command here to illustrate its usefulness and ease of use.

- [Merge Command](#merge-command)

## Merge Command

Use the `merge` command to merge several images from different directories into one.

```
Usage:
  generator merge [flags]

Flags:
  -d, --directory strings   specify a directory       
  -h, --help                help for merge
  -o, --out string          export the images
  -D, --root string         specify the root directory
```

let me try to explain how this tool work, please refer to the following diagram:

```
@startuml Basic Sample
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

HIDE_STEREOTYPE()

Boundary(c1, "Background Folder") {
  Container(bg_images, "Images", "png files", "Specify the background images for merging\n\n example: blue, red, yellow, etc.")
}

Boundary(c2, "Body Folder") {
  Container(body_images, "Images", "png files", "Specify the body images for merging\n\n example: different skin colors")
}

Boundary(c3, "Face Folder") {
  Container(face_images, "Images", "png files", "Specify the face images for merging\n\n example: the different combina - eye and mouth")
}

Boundary(c4, "Output Folder") {
  Container(com, "Image combinations", "png files")
  Container(file, "result.json", "json file", "The results file provides information about the source of image combinations.")
}


System(command, "Merge Command")

Rel(bg_images, command, "add")
Rel(body_images, command, "add")
Rel(face_images, command, "add")
Rel(command, com, "export")
Rel(command, file, "export")
@enduml
```

