# pngfy

**Author:** Pietro Mascolo

pngfy converts a directory of pdf files into individual png files for each page


Usage for single file:
```bash
./pngfy FILE [--width=WIDTH] [--height=HEIGHT][--target=TARGET]
```

Usage for whole directory:
```bash
./pngfy convert DIRECTORY [--width=WIDTH] [--height=HEIGHT][--target=TARGET]
```

Width and height are optional parameters to set the size of the resulting images (in pixels). Default values are 840x1188 (A4 ratio).  
Target represents the parent folder to the result directories, default target is `<DIRECTORY>/target` or `<directory of FILE>/target`.

The results will be saved in target directory: one sub-directory per file  containing png images numbered from `0` to `N`.