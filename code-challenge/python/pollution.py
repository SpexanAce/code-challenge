import json
from datetime import datetime

def parse_pollution_data(data):
    morning_pm10 = 0
    morning_pm25 = 0
    afternoon_pm10 = 0
    afternoon_pm25 = 0
    night_pm10 = 0
    night_pm25 = 0

    for idx, timeOfDay in enumerate(data['hourly']['time']):
        dt = datetime.strptime(timeOfDay, '%Y-%m-%dT%H:%M')
        
        if dt.hour > 6 and dt.hour < 12:
            morning_pm10 += data['hourly']['pm10'][idx]
            morning_pm25 += data['hourly']['pm2_5'][idx]
            continue
        
        if dt.hour > 12 and dt.hour < 18:
            afternoon_pm10 += data['hourly']['pm10'][idx]
            afternoon_pm25 += data['hourly']['pm2_5'][idx]
            continue

        if dt.hour < 6 or dt.hour > 18:
            night_pm10 += data['hourly']['pm10'][idx]
            night_pm25 += data['hourly']['pm2_5'][idx]
            continue
    
    return {
        'pm10': {
            'morning': morning_pm10 / 6,
            'afternoon': afternoon_pm10 / 6,
            'night': night_pm10 / 12,
        },
        'pm25': {
            'morning': morning_pm25 / 6,
            'afternoon': afternoon_pm25 / 6,
            'night': night_pm25 / 12,
        }
    }

def load_pollution_data():
    fh = open('../database/pollution_data.json')
    data = json.load(fh)
    fh.close()

    return data

def print_pollution_data():
    raw = load_pollution_data()
    data = parse_pollution_data(raw)

    print(f"PM-10 averages:")
    print(f"\tMorning: {data['pm10']['morning']:.2f} μg/m³")
    print(f"\tAfternoon: {data['pm10']['afternoon']:.2f} μg/m³")
    print(f"\tNight: {data['pm10']['night']:.2f} μg/m³")
    
    print(f"PM-2.5 averages:")
    print(f"\tMorning: {data['pm25']['morning']:.2f} μg/m³")
    print(f"\tAfternoon: {data['pm25']['afternoon']:.2f} μg/m³")
    print(f"\tNight: {data['pm25']['night']:.2f} μg/m³")

if __name__ == '__main__':
    print_pollution_data()