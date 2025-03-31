using System;
using System.IO;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace PollutionData;

class Program
{
    static async Task Main(string[] args)
    {
        try
        {
            var data = await LoadPollutionDataAsync();
            var da = ParsePollutionData(data);

            Console.WriteLine("PM-10 averages:");
            Console.WriteLine($"\tMorning: {da.Pm10.Morning:0.00} μg/m³");
            Console.WriteLine($"\tAfternoon: {da.Pm10.Afternoon:0.00} μg/m³");
            Console.WriteLine($"\tNight: {da.Pm10.Night:0.00} μg/m³");

            Console.WriteLine("PM-2.5 averages:");
            Console.WriteLine($"\tMorning: {da.Pm25.Morning:0.00} μg/m³");
            Console.WriteLine($"\tAfternoon: {da.Pm25.Afternoon:0.00} μg/m³");
            Console.WriteLine($"\tNight: {da.Pm25.Night:0.00} μg/m³");
        }
        catch (Exception ex)
        {
            Console.WriteLine($"An error occurred: {ex.Message}");
        }
    }

    private static async Task<byte[]> LoadPollutionDataAsync()
    {
        return await File.ReadAllBytesAsync("../database/pollution_data.json");
    }

    private static DailyAverages ParsePollutionData(byte[] rawData)
    {
        var data = JsonSerializer.Deserialize<HourlyData>(rawData);

        double morningPm10 = 0, morningPm25 = 0;
        double afternoonPm10 = 0, afternoonPm25 = 0;
        double nightPm10 = 0, nightPm25 = 0;

        for (int i = 0; i < data.Hourly.Time.Length; i++)
        {
            var time = DateTime.Parse(data.Hourly.Time[i]);

            if (time.Hour > 6 && time.Hour < 12)
            {
                morningPm10 += data.Hourly.Pm10[i];
                morningPm25 += data.Hourly.Pm25[i];
            }
            else if (time.Hour > 12 && time.Hour < 18)
            {
                afternoonPm10 += data.Hourly.Pm10[i];
                afternoonPm25 += data.Hourly.Pm25[i];
            }
            else if (time.Hour < 6 || time.Hour > 18)
            {
                nightPm10 += data.Hourly.Pm10[i];
                nightPm25 += data.Hourly.Pm25[i];
            }
        }

        var da = new DailyAverages
        {
            Pm10 = new Averages
            {
                Morning = morningPm10 / 6,
                Afternoon = afternoonPm10 / 6,
                Night = nightPm10 / 12
            },
            Pm25 = new Averages
            {
                Morning = morningPm25 / 6,
                Afternoon = afternoonPm25 / 6,
                Night = nightPm25 / 12
            }
        };

        return da;
    }
}

public class DailyAverages
{
    public Averages Pm10 { get; set; }
    public Averages Pm25 { get; set; }
}

public class Averages
{
    public double Morning { get; set; }
    public double Afternoon { get; set; }
    public double Night { get; set; }
}

public class HourlyData
{
    [JsonPropertyName("hourly")]
    public Hourly Hourly { get; set; }
}

public class Hourly
{
    [JsonPropertyName("time")]
    public string[] Time { get; set; }
    
    [JsonPropertyName("pm10")]
    public double[] Pm10 { get; set; }
    
    [JsonPropertyName("pm2_5")]
    public double[] Pm25 { get; set; }
}