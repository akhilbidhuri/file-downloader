# File Downloader

This is a concurrent file downloader made using GO

## Usage

```sh
git clone https://github.com/akhilbidhuri/file-downloader.git
cd file-downloader
go run main.go <url> <out_put_file_path>
```
if out_put_file_path is not provided a new file will be created in the current directory deriving its name from url 

### or
build the project-
```sh
go build -o bin/downloader cmd/main.go
cd bin
./downloader <url> [output_file_path]
```
## Example

```sh
go run .\cmd\main.go https://svs.gsfc.nasa.gov/vis/a030000/a030800/a030877/frames/5760x3240_16x9_01p/BlackMarble_2016_1200m_africa_s_labeled.png
```
