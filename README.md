# msds434_final

[Video Walkthrough](https://northwestern.hosted.panopto.com/Panopto/Pages/Viewer.aspx?id=648d003b-a024-4b85-b5d0-b0cd01851733&start=0)

## 1. Local Machine
    -Pull data from: https://www.transtats.bts.gov/OT_Delay/OT_DelayCause1.asp?20=E
    -We pulled 10 years (2010-2020)
    -Run the data_prep.ipynb file, which creates the following:
        Two testing and training sets (four each)
        Carrier Name and Airport Code keys - SageMaker only accepts numbers
        File paths will need to be adjusted
        

## 2. AWS Steps
    -Create an IAM role with the following permissions (unsure if all are still needed):
        AdministratorAccess-AWSElasticBeanstalk
        AmazonEC2ContainerRegistryFullAccess
        AmazonEC2FullAccess
        AmazonS3FullAccess
        AmazonSageMakerFullAccess
        AWSElasticBeanstalkWebTier
        AWSElasticBeanstalkWorkerTier
        CloudWatchFullAccess
        CloudWatchFullAccessV2
        IAMFullAccess
    -Start an Elastic Beanstalk application. This will use the role created above and should have the following settings: 
        Docker running on 64bit Amazon Linux 2
        Everything else can be defaulted (including the code).
        Navigate to EC2 instance and select the newly created instace. Update the persmissions to allow for all incoming traffic.
    -Start a SageMaker Studio session (default everything including user profile information)
        Upload the testing and training data to the automatically created S3 buckets. Record the bucket name. 
    -Create two Jupyter notebooks in SageMaker 
        Change the instance to ml.m5.large
        Copy the following code over:
            model_1_time.py
            model_2_odds.py
        Update the bucket_name variable for both from 'tktk' to the bucket created above
        Run both models and record their endpoints

## 3. github
    -You must upload the following to the repositories' secrets keys:
        AWS_ACCESS_KEY_ID
        AWS_SECRET_ACCESS_KEY
    -Your github must have the following structure:
        .github
            workflows
                deploy.yml
        data
            airport_classes.json
            carrier_classes.json
        Dockerfile
        Dockerrun.aws.json
        go.mod
        go.sum
        main.go
    -For deploy.yml, the following variables must match your AWS Beanstalk application
        application_name
        enviroment_name
    -main.go: Update the endpoints with those you recorded above. The should look like this:
        var endpoint1 = "linear-learner-2023-12-02-18-11-46-657"
        var endpoint2 = "linear-learner-2023-12-02-17-57-03-454"
    -Commit changes:
        This will update the code and dockerize it as outlined in the deploy.yml and Dockerfile. 
        You should be able to see the updated version on the Elastic Beanstalk dashboard
            CloudWatch Metrics are available automatically through Elastic Beanstalk and Cloudwatch

## 4. Test
    -Testing can occur with any of of the relevant months (numeral), carriers (string) and airport code (string). Sample below:
    curl -X POST -H "Content-Type: application/json" -d '{"month":12, "carrier":"American Airlines Inc.", "airport":"PHX"}' http://airlines-env.eba-geqwdvnt.us-east-2.elasticbeanstalk.com/predict
