This folder contains the apis that the frontend hits for info. These will all have some kind of parser that parses the raw API data, then 
returns the parsed API data for the frontend to later use. 

The types and apis that are avaliable to the frontend are summarized in the `./spec` folder. 

The general structure followed here is 

`api.go` will have the external API function that is the only entrypoint into this API. 

`implementations.go` will hold the implementations for the `InternalAPI` methods. 

`types.go` will hold the external API type that we parse from our GET request. 

`utils.go` will hold the util functions. 

If you think something can be made generic, consider putting it in the higher level `utils/` package. 
