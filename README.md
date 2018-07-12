# XMen REST API

REST API that detects whether a dna corresponds to a human or a mutant

| resource      | description                       |
|:--------------|:----------------------------------|
| `/mutant`      | detects if a human is mutant by means of the DNA sequence
| `/stats`    | returns a Json with the statistics from the DNA verifications

In order to follow these steps, you'll need to have [Docker](https://www.docker.com/) installed.

## Build the image
```bash
~$ docker build -t xmen .
```
## Run the container

Set up a Docker container to run the api.

```bash
~$ docker run -p 8080:8080 -d xmen
```

## Detects if a human is mutant

Now, lets get some data.

```bash
~$ curl --header "Content-Type: application/json" \
     --request POST \
     --data '{
   	"dna":["ATGCGA","CCGTGC","TTATGT","AGAAGG","CCACTA","TCACTG"]}' \
     http://localhost:8080/mutant
```
Will return the following result:

```json
{
    "isMutant": true
}
```

## Obtain stats


```bash
~$ curl http://localhost:8080/stats
```

