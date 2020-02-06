This folder contains the types that will be sent to the frontend. They are _not_ the types which are recieved 
from the raw APIs - those should be kept in their respective parser files. 

The type files should also contain the implemented `Serialize()` function for that specific type (so that it 
implements the `InternalAPI` interface)
