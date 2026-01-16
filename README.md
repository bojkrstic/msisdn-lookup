# msisdn-lookup

For local development, you have two options:

Via Docker:

1. Build the image once: docker build -t msisdn-lookup .

2. Run it whenever you need it: docker run --rm -p 8080:8080 --name msisdn msisdn-lookup
or this -> docker run --rm -p 9090:9090 --name msisdn msisdn-lookup

3. Visit http://localhost:8080 (or set the port as desired, e.g. -p 9090:8080).

4. For a quick reload after code changes, run the build and docker run again. If you want "watch" behavior, you can use docker buildx bake --load or similar scripts, but the simplest is to rebuild after major changes.


Rebild on server:


1. cd /opt/msisdn-lookup
2. git pull
3. docker build -t msisdn-lookup:latest .
4. sudo systemctl restart msisdn-lookup
5. http://83-229-82-132.cloud-xip.com/msisdn/
