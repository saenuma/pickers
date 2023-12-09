# pickers

a file and color picker for embedding as an executable.

## Why this project

Other dialogs give problems while embedding in snapcraft. zenity for instance increases 
the size of the build by about 120MB.

This project provides a static build.

## file picker
The file picker is called fpicker. It is quite basic

It expects a path and extensions.

For example `fpicker /home/b/snap/flaar/common txt|bin|blob` 

## Color Picker
The color picker is called cpicker. It is quite basic

It doesn't expect any argument.
