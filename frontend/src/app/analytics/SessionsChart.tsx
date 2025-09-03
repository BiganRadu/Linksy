import * as React from 'react';
import { useTheme } from '@mui/material/styles';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import Stack from '@mui/material/Stack';
import { LineChart } from '@mui/x-charts/LineChart';
import Cookies from 'js-cookie';
import axios from 'axios';

function AreaGradient({ color, id }: { color: string; id: string }) {
  return (
    <defs>
      <linearGradient id={id} x1="50%" y1="0%" x2="50%" y2="100%">
        <stop offset="0%" stopColor={color} stopOpacity={0.5} />
        <stop offset="100%" stopColor={color} stopOpacity={0} />
      </linearGradient>
    </defs>
  );
}


function getNameOfTheDays(startTimestamp: number, endTimestamp: number) : string[] {
	const days = [];
	const today = new Date(endTimestamp * 1000);
	let numOfDays = Math.floor((endTimestamp - startTimestamp) / (60 * 60 * 24));
	for (let i = numOfDays; i >= 1; i--) {
		const d = new Date(today);
		d.setDate(today.getDate() - i);
		const monthName = d.toLocaleDateString('en-US', { month: 'short' });
		days.push(`${monthName} ${d.getDate()}`);
	}
	return days;
}

export default function SessionsChart({startTimestamp, endTimestamp}: {startTimestamp: number, endTimestamp: number}) {
  const theme = useTheme();
  const [days, setDays] = React.useState<string[]>(getNameOfTheDays(startTimestamp, endTimestamp));
  const [data, setData] = React.useState(Array(days.length).fill(0));
  const [totalSessions, setTotalSessions] = React.useState(0);

  	React.useEffect(() => {
	const fetchData = async () => {
		console.log(`Fetching data from ${startTimestamp} to ${endTimestamp}`);
  
		const authToken = Cookies.get('AuthToken');
    axios.get(`https://linksy-mhe5.onrender.com/app/analytics?chart_code=sessions&start=${startTimestamp}&end=${endTimestamp}`, {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			setData(response.data.sessions);
			setDays(getNameOfTheDays(startTimestamp, endTimestamp));
			setTotalSessions(response.data.total);
		}).catch(error => {
			console.log("big problem");
		});
	};
  
	fetchData();
	  }, [startTimestamp, endTimestamp]);

  const colorPalette = [
    theme.palette.primary.light,
    theme.palette.primary.main,
    theme.palette.primary.dark,
  ];

  return (
    <Card variant="outlined" sx={{ width: '100%' }}>
      <CardContent>
        <Typography component="h2" variant="subtitle2" gutterBottom>
          Accesses
        </Typography>
        <Stack sx={{ justifyContent: 'space-between' }}>
          <Stack
            direction="row"
            sx={{
              alignContent: { xs: 'center', sm: 'flex-start' },
              alignItems: 'center',
              gap: 1,
            }}
          >
            <Typography variant="h4" component="p">
              {totalSessions}
            </Typography>
          </Stack>
          <Typography variant="caption" sx={{ color: 'text.secondary' }}>
            Accesses per day for the selected period
          </Typography>
        </Stack>
        <LineChart
          colors={colorPalette}
          xAxis={[
            {
              scaleType: 'point',
              data: days,
              tickInterval: (index, i) => (i + 1) % 5 === 0,
              height: 24,
            },
          ]}
          yAxis={[{ width: 50 }]}
          series={[
            {
              id: 'direct',
              label: '',
              showMark: false,
              curve: 'linear',
              stack: 'total',
              area: true,
              stackOrder: 'ascending',
              data: data,
            },
          ]}
          height={250}
          margin={{ left: 0, right: 20, top: 20, bottom: 0 }}
          grid={{ horizontal: true }}
          sx={{
            '& .MuiAreaElement-series-direct': {
              fill: "url('#direct')",
            },
          }}
          hideLegend
        >
          <AreaGradient color={theme.palette.primary.light} id="direct" />
        </LineChart>
      </CardContent>
    </Card>
  );
}
