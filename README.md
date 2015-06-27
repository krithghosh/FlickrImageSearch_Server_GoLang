# Flickr Api Call and Database Interaction - Web Server

The Go App Server is responsible for getting the response from the flickr apis and interact with the database.

The flickr api call : https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=6b54d86b4e09671ef6a2a8c02b7a3537&text=cute+puppies&format=json&nojsoncallback=1

On getting the response from the flickr it stores the json data in the MongoLab. The flickr api sends in data about the image like server id, id, secret id, farm_id and these values are used to form the url for the pictures : https://farm{farm-id}.staticflickr.com/{server-id}/{id}_{secret}.jpg

When this server is called using /UpVote, then the link of the image is extracted from the query string and used to increases the upvote of the image in the MongoLab.

When this server is called using /DownVote, then the link of the image is extracted from the query string and used to increases the downvote of the image in the MongoLab.
