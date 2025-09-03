import * as React from 'react';
import { PieChart } from '@mui/x-charts/PieChart';
import { useDrawingArea } from '@mui/x-charts/hooks';
import { styled } from '@mui/material/styles';
import Typography from '@mui/material/Typography';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import LinearProgress, { linearProgressClasses } from '@mui/material/LinearProgress';
import axios from 'axios';
import Cookies from 'js-cookie';


interface StyledTextProps {
  variant: 'primary' | 'secondary';
}

const StyledText = styled('text', {
  shouldForwardProp: (prop) => prop !== 'variant',
})<StyledTextProps>(({ theme }) => ({
  textAnchor: 'middle',
  dominantBaseline: 'central',
  fill: (theme.vars || theme).palette.text.secondary,
  variants: [
    {
      props: {
        variant: 'primary',
      },
      style: {
        fontSize: theme.typography.h5.fontSize,
      },
    },
    {
      props: ({ variant }) => variant !== 'primary',
      style: {
        fontSize: theme.typography.body2.fontSize,
      },
    },
    {
      props: {
        variant: 'primary',
      },
      style: {
        fontWeight: theme.typography.h5.fontWeight,
      },
    },
    {
      props: ({ variant }) => variant !== 'primary',
      style: {
        fontWeight: theme.typography.body2.fontWeight,
      },
    },
  ],
}));

interface PieCenterLabelProps {
  primaryText: string;
  secondaryText: string;
}

function PieCenterLabel({ primaryText, secondaryText }: PieCenterLabelProps) {
  const { width, height, left, top } = useDrawingArea();
  const primaryY = top + height / 2 - 10;
  const secondaryY = primaryY + 24;

  return (
    <React.Fragment>
      <StyledText variant="primary" x={left + width / 2} y={primaryY}>
        {primaryText}
      </StyledText>
      <StyledText variant="secondary" x={left + width / 2} y={secondaryY}>
        {secondaryText}
      </StyledText>
    </React.Fragment>
  );
}

export default function ChartUserByCountry({chart_code, startTimestamp, endTimestamp}: {chart_code: string, startTimestamp: number, endTimestamp: number}) {
	const [pairs, setPairs] = React.useState([]);
	const [total, setTotal] = React.useState(0);
	let colors = [];
	let data = [];
	React.useEffect(() => {
	const fetchData = async () => {
		const authToken = Cookies.get('AuthToken');
    axios.get(`https://linksy-mhe5.onrender.com/app/analytics?chart_code=${chart_code}&start=${startTimestamp}&end=${endTimestamp}`, {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			setPairs(response.data.values);
			setTotal(response.data.total);
		}).catch(error => {
			console.log("Error fetching analytics data:", error);
		});
	};
  
	fetchData();
	  }, [startTimestamp, endTimestamp]);
	for (let i = 0; i < pairs.length; i++) {
		data.push({
			name: pairs[i].name,
			value: pairs[i].value,
			color: `hsl(220, 20%, ${65 - 40 / (pairs.length - i)}%)`,
		});
		colors.push(`hsl(220, 20%, ${65 - 40 / (pairs.length - i)}%)`);
	}

	if (data.length === 0) {
		data.push({
			name: 'No data',
			value: 1,
			color: 'hsl(220, 20%, 65%)',
		});
		colors.push('hsl(220, 20%, 65%)');
	}
  return (
    <Card
      variant="outlined"
      sx={{ display: 'flex', flexDirection: 'column', gap: '8px', flexGrow: 1 }}
    >
      <CardContent>
        <Typography component="h2" variant="subtitle2">
          Users by country
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <PieChart
            colors={colors}
            margin={{
              left: 80,
              right: 80,
              top: 80,
              bottom: 80,
            }}
            series={[
              {
                data,
                innerRadius: 75,
                outerRadius: 100,
                paddingAngle: 0,
                highlightScope: { fade: 'global', highlight: 'item' },
              },
            ]}
            height={260}
            width={260}
            hideLegend
          >
            <PieCenterLabel primaryText={total} secondaryText="Total" />
          </PieChart>
        </Box>
        {data.map((pair, index) => (
          <Stack
            key={index}
            direction="row"
            sx={{ alignItems: 'center', gap: 2, pb: 2 }}
          >
            <Stack sx={{ gap: 1, flexGrow: 1 }}>
              <Stack
                direction="row"
                sx={{
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  gap: 2,
                }}
              >
                <Typography variant="body2" sx={{ fontWeight: '500' }}>
                  {pair.name}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  {pair.value/total * 100}%
                </Typography>
              </Stack>
              <LinearProgress
                variant="determinate"
                aria-label="Number of users by country"
                value={pair.value/total * 100}
                sx={{
                  [`& .${linearProgressClasses.bar}`]: {
                    backgroundColor: pair.color,
                  },
                }}
              />
            </Stack>
          </Stack>
        ))}
      </CardContent>
    </Card>
  );
}
