import numpy as np
from sklearn.preprocessing import LabelEncoder
import pandas as pd
from sklearn.model_selection import train_test_split
import joblib
import json

# directory path
desktop_directory = '/Users/zacharywatson/Desktop/434Final/Week10/'
df = pd.read_csv('/Users/zacharywatson/Desktop/434Final/Week10/Airline_Delay_Cause.csv')

# Fill null values with zero
df = df.fillna(0)

# Select the relevant columns
relevant_columns = ['arr_del15', 'month', 'carrier_name', 'airport', 'arr_flights', 'arr_delay']
df = df[relevant_columns]

airlines = ['American Airlines Inc.', 'Alaska Airlines Inc.', 'JetBlue Airways', 'Delta Air Lines Inc.', 'Frontier Airlines Inc.', 'Allegiant Air', 'Hawaiian Airlines Inc.', 'Spirit Air Lines', 'United Air Lines Inc.', 'Southwest Airlines Co.']

df = df[df['carrier_name'].isin(airlines)]

le = LabelEncoder()

df['carrier_name'] = le.fit_transform(df['carrier_name'])
carrier_classes = {cls: i for i, cls in enumerate(le.classes_)}
with open(f'{desktop_directory}Upload/carrier_classes.json', 'w') as f:
    json.dump(carrier_classes, f)

df['airport'] = le.fit_transform(df['airport'])
airport_classes = {cls: i for i, cls in enumerate(le.classes_)}
with open(f'{desktop_directory}Upload/airport_classes.json', 'w') as f:
    json.dump(airport_classes, f)

# Create a copy of the dataframe for the second model
df2 = df.copy()

# Calculate the ratio of 'arr_delay' to 'arr_del15' for the original dataframe
df['arr_delay'] = np.where(df['arr_del15'] == 0, 0, df['arr_delay'] / df['arr_del15'])

# Drop the 'arr_del15' and 'arr_flights' columns from the original dataframe
df = df.drop(['arr_del15', 'arr_flights'], axis=1)

# Rearrange the columns so that 'arr_delay' is the first column
cols1 = ['arr_delay']  + [col for col in df if col != 'arr_delay']
df = df[cols1]

# Calculate the odds of delay for the second dataframe
df2['odds_of_delay'] = np.where(df2['arr_flights'] == 0, 0, df2['arr_del15'] / df2['arr_flights'])

# Drop the 'arr_delay', 'arr_del15' and 'arr_flights' columns from the second dataframe
df2 = df2.drop(['arr_delay','arr_del15', 'arr_flights'], axis=1)

# Rearrange the columns so that 'odds_of_delay' is the first column
cols2 = ['odds_of_delay']  + [col for col in df2 if col != 'odds_of_delay']
df2 = df2[cols2]

# Split the first dataframe into training and testing data
train_df1, test_df1 = train_test_split(df, test_size=0.2)

# Split the second dataframe into training and testing data
train_df2, test_df2 = train_test_split(df2, test_size=0.2)

# Convert the DataFrame to a CSV file 
train_df1.to_csv(f'{desktop_directory}AWS/TrainTest/train1.csv', index=False, header=False)
test_df1.to_csv(f'{desktop_directory}AWS/TrainTest/test1.csv', index=False, header=False)

train_df2.to_csv(f'{desktop_directory}AWS/TrainTest/train2.csv', index=False, header=False)
test_df2.to_csv(f'{desktop_directory}AWS/TrainTest/test2.csv', index=False, header=False)
