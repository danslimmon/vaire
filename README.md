# Vairë

Docker image manager and deployer for use with OpenShift.

## Examples

### See all images in the `qaready` queue

```
$ curl -H 'Token AbCdEf...' -XPOST --data '{"image_id":"3b0d1a926a", "app": "some-application", "comment": "Submitted by Jane Doe"}' https://vaire.example.com/api/v1/queues
```

```
$ curl -H 'Token AbCdEf...' https://vaire.example.com/api/v1/queues/qaready
{
  "error": null,
  "result": [
    {
      "image_id": "3b0d1a926a",
      "app": "some-application",
      "time": "2016-09-29T18:46:13+0000"
    }
  ]
```
