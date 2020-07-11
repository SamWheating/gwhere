# gwhere

As in, _"Gee, where is that bucket anyways?"_

Also as in, _"Google, where did I put that bucket?"_

Mini Cli tool to remind you which GCP project a storage bucket is in.

## Usage

`gwhere <bucket>`

You'll need `stat` access to the bucket and access to its containg project.

## Why does this exist?

Previously it took me two `gsutil` / `gcloud` commands to get the project ID of one of my buckets. I had to do this on at least two separate occasions this year so it was becoming a significant time sink. 

## Couldn't this just be done with Bash / gcloud / jq / etc?

Probably, yeah. 
