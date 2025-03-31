import fs from 'fs';

interface IPollutionData {
  latitude: number
  longitude: number
  generationtime_ms: number
  utc_offset_seconds: number
  timezone: string
  timezone_abbreviation: string
  hourly_units: {
    time: string
    pm10: string
    pm2_5:string
  }
  hourly: {
    time: string[]
    pm10: number[]
    pm2_5: number[]
  }
}


const parsePollutionData = (data: IPollutionData) => {
    let morningPm10 = 0;
    let morningPm25 = 0;
    let afternoonPm10 = 0;
    let afternoonPm25 = 0;
    let nightPm10 = 0;
    let nightPm25 = 0;

    data.hourly.time.forEach((timestamp, idx) => {
      const hour = new Date(timestamp).getHours();

      if (hour > 6 && hour < 12) {
        morningPm10 += data.hourly.pm10[idx];
        morningPm25 += data.hourly.pm2_5[idx];
      }

      if (hour > 12 && hour < 18) {
        afternoonPm10 += data.hourly.pm10[idx];
        afternoonPm25 += data.hourly.pm2_5[idx];
      }

      if (hour < 6 || hour > 18) {
        nightPm10 += data.hourly.pm10[idx];
        nightPm25 += data.hourly.pm2_5[idx];
      }
    });
    
    return {
        pm10: {
            morning: morningPm10 / 6,
            afternoon: afternoonPm10 / 6,
            night: nightPm10 / 12,
        },
        pm25: {
            morning: morningPm25 / 6,
            afternoon: afternoonPm25 / 6,
            night: nightPm25 / 12,
        }
    }
  }

const loadPollutionData = () =>
  JSON.parse(fs.readFileSync('../database/pollution_data.json').toString()) as IPollutionData

const printPollutionData = () => {
  const raw = loadPollutionData();
  const data = parsePollutionData(raw);

  console.log('PM-10 averages:')
  console.log(`\tMorning: ${data.pm10.morning} μg/m³`)
  console.log(`\tAfternoon: ${data.pm10.afternoon} μg/m³`);
  console.log(`\tNight: ${data.pm10.night} μg/m³`);
    
  console.log('PM-2.5 averages:');
  console.log(`\tMorning: ${data.pm25.morning} μg/m³`);
  console.log(`\tAfternoon: ${data.pm25.afternoon} μg/m³`);
  console.log(`\tNight: ${data.pm25.night} μg/m³`);
}

printPollutionData();
