curl -s "http://localhost:3000/?format=json&to=2023-02-01" | jq > ./testdata/ResponseFilter.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&orderby=modified" | jq > ./testdata/ResponseFilter-modified.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&orderby=modified&direction=asc" | jq > ./testdata/ResponseFilter-modified-asc.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&term=copp" | jq > ./testdata/ResponseFilter-term.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&term=su+špa" | jq > ./testdata/ResponseFilter-term-1.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&camera=NIKON%20D800" | jq > ./testdata/ResponseFilter-camera.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&lens=AF-S%20Nikkor%2050mm%20f%2f1.8G" | jq > ./testdata/ResponseFilter-lens.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&lens=123" | jq > ./testdata/ResponseFilter-lens-1.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&lens=Nikon%20Ai-s%20105mm%20f%2f2.5" | jq > ./testdata/ResponseFilter-lens-2.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&folder=2017/20170624-idaho" | jq > ./testdata/ResponseFilter-folder.json
curl -s "http://localhost:3000/?format=json&to=2023-02-01&tag=idaho" | jq > ./testdata/ResponseFilter-tag.json
curl -s "http://localhost:3000/onthisday?format=json" | jq > ./testdata/ResponseFilter-onthisday.json

curl -s "http://localhost:3000/media/651935749?format=json" | jq > ./testdata/ResponseImage-0.json
curl -s "http://localhost:3000/media/3455659031?format=json" | jq > ./testdata/ResponseImage-1.json
curl -s "http://localhost:3000/media/4119775194?format=json" | jq > ./testdata/ResponseImage-2.json
curl -s "http://localhost:3000/media/651935749/in/folder/2017/20170624-idaho?format=json" | jq > ./testdata/ResponseImage-folder.json
curl -s "http://localhost:3000/media/651935749/in/tag/idaho?format=json" | jq > ./testdata/ResponseImage-tag.json
curl -s "http://localhost:3000/media/4119775194/in/tag/%40acconfb?format=json" | jq > ./testdata/ResponseImage-tag-acc.json
curl -s "http://localhost:3000/media/4119775194/in/tag/%23californiawildfires?format=json" | jq > ./testdata/ResponseImage-tag-cal.json
curl -s "http://localhost:3000/media/651935749/in/favorites?format=json" | jq > ./testdata/ResponseImage-favorites.json

# prev/next responses
curl -s "http://localhost:3000/media/3455659031?camera=NIKON%20D800&format=json" | jq > ./testdata/ResponseImage-camera.json
curl -s "http://localhost:3000/media/651935749?lens=AF-S%20Nikkor%2050mm%20f%2f1.8G&format=json" | jq > ./testdata/ResponseImage-lens.json
curl -s "http://localhost:3000/media/525791494?lens=Nikon%20Ai-s%20105mm%20f%2f2.5&format=json" | jq > ./testdata/ResponseImage-lens-1.json
curl -s "http://localhost:3000/media/264898052?focallength35=50&format=json" | jq > ./testdata/ResponseImage-focallength35.json
curl -s "http://localhost:3000/media/3216513272?software=darktable%204.4.2&format=json" | jq > ./testdata/ResponseImage-software.json
curl -s "http://localhost:3000/media/1129346697?term=bogus&format=json" | jq > ./testdata/ResponseImage-term.json

curl -s "http://localhost:3000/favorites?format=json" | jq > ./testdata/ResponseImages-favorites.json

curl -s "http://localhost:3000/folders?format=json" | jq > ./testdata/ResponseFolders.json
curl -s "http://localhost:3000/folder/2019/20190330-sawtooths?format=json" | jq > ./testdata/ResponseFolder.json

curl -s "http://localhost:3000/tags?format=json" | jq > ./testdata/ResponseTags.json
curl -s "http://localhost:3000/tag/idaho?format=json" | jq > ./testdata/ResponseTag.json
curl -s "http://localhost:3000/tag/%40acconfb?format=json" | jq > ./testdata/ResponseTag-acc.json
curl -s "http://localhost:3000/tag/%23californiawildfires?format=json" | jq > ./testdata/ResponseTag-cal.json

curl -s "http://localhost:3000/gear?format=json" | jq > ./testdata/ResponseGear.json

curl -s "http://localhost:3000/map?format=json" | jq > ./testdata/ResponseMap.json
