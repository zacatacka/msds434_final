import boto3
import sagemaker
import numpy as np
from sagemaker import get_execution_role
from sagemaker.inputs import TrainingInput
from sagemaker.image_uris import retrieve

# bucket name and data paths
bucket_name = 'sagemaker-us-east-2-777338306685'
train_data_path = 's3://{}/train2.csv'.format(bucket_name)
test_data_path = 's3://{}/test2.csv'.format(bucket_name)

# URI of the Linear Learner algorithm in your region
container = retrieve('linear-learner', sagemaker.Session().boto_region_name)

# sageMaker execution role
role = get_execution_role()

s3_input_train = TrainingInput(s3_data=train_data_path, content_type='text/csv')
s3_input_test = TrainingInput(s3_data=test_data_path, content_type='text/csv')

estimator = sagemaker.estimator.Estimator(container,
                                          role, 
                                          instance_count=1, 
                                          instance_type='ml.m5.large',
                                          output_path='s3://{}/{}/output'.format(bucket_name, 'linear-learner'),
                                          sagemaker_session=sagemaker.Session(),
                                          enable_cloudwatch_metrics=True)

# hyperparameters
estimator.set_hyperparameters(predictor_type='binary_classifier')

# fit
estimator.fit({'train': s3_input_train, 'validation': s3_input_test}, logs=False)

# deploy
predictor = estimator.deploy(initial_instance_count=1, instance_type='ml.m5.large')

# print endpoint name
print(predictor.endpoint_name)

# define sigmoid function
#def sigmoid(x):
#    return 1 / (1 + np.exp(-x))

# get predictions 
#predictions = predictor.predict(test_data)

# apply sigmoid 
#probabilities = sigmoid(predictions)
