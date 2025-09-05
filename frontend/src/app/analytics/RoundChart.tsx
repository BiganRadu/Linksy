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
    { props: { variant: 'primary' }, style: { fontSize: theme.typography.h5.fontSize, fontWeight: theme.typography.h5.fontWeight } },
    { props: ({ variant }) => variant !== 'primary', style: { fontSize: theme.typography.body2.fontSize, fontWeight: theme.typography.body2.fontWeight } },
  ],
}));

interface PieCenterLabelProps {
  primaryText: string | number;
  secondaryText: string;
}

function PieCenterLabel({ primaryText, secondaryText }: PieCenterLabelProps) {
  const { width, height, left, top } = useDrawingArea();
  const primaryY = top + height / 2 - 10;
  const secondaryY = primaryY + 24;

  return (
    <>
      <StyledText variant="primary" x={left + width / 2} y={primaryY}>
        {primaryText}
      </StyledText>
      <StyledText variant="secondary" x={left + width / 2} y={secondaryY}>
        {secondaryText}
      </StyledText>
    </>
  );
}

export default function ChartUserByCountry({
  chart_code,
  startTimestamp,
  endTimestamp,
}: {
  chart_code: string;
  startTimestamp: number;
  endTimestamp: number;
}) {
  const [pairs, setPairs] = React.useState<Array<{ name: string; value: number }>>([]);
  const [total, setTotal] = React.useState(0);

  React.useEffect(() => {
    const fetchData = async () => {
      const authToken = Cookies.get('AuthToken');
      axios
        .get(
          `https://linksy-mhe5.onrender.com/app/analytics?chart_code=${chart_code}&start=${startTimestamp}&end=${endTimestamp}`,
          {
            headers: { AuthToken: authToken },
          },
        )
        .then((response) => {
          setPairs(response.data.values ?? []);
          setTotal(response.data.total ?? 0);
        })
        .catch((error) => {
          console.log('Error fetching analytics data:', error);
        });
    };

    fetchData();
  }, [chart_code, startTimestamp, endTimestamp]);

  // Build colors and chart data with labels
  const colors: string[] = [];
  const dataChart: Array<{ id: number; value: number; label: string; color: string }> = [];

  for (let i = 0; i < pairs.length; i++) {
    const color = `hsl(220, 20%, ${65 - 40 / (pairs.length - i)}%)`;
    colors.push(color);
    dataChart.push({
      id: i,
      value: pairs[i].value,
      label: pairs[i].name, // label is used by the tooltip
      color,
    });
  }

  // Fallback slice so the chart doesn’t look empty
  if (dataChart.length === 0) {
    dataChart.push({
      id: 0,
      value: 1,
      label: 'No data',
      color: 'hsl(220, 20%, 65%)',
    });
    colors.push('hsl(220, 20%, 65%)');
  }

  return (
    <Card variant="outlined" sx={{ display: 'flex', flexDirection: 'column', gap: '8px', flexGrow: 1 }}>
      <CardContent>
        <Typography component="h2" variant="subtitle2">
          Users by country
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <PieChart
            colors={colors}
            margin={{ left: 80, right: 80, top: 80, bottom: 80 }}
            series={[
              {
                data: dataChart,
                innerRadius: 75,
                outerRadius: 100,
                paddingAngle: 0,
                highlightScope: { fade: 'global', highlight: 'item' },
                valueFormatter: (item) => {
                  return `${item.value}`;
                },
              },
            ]}
            height={260}
            width={260}
            hideLegend
          >
            <PieCenterLabel primaryText={total} secondaryText="Total" />
          </PieChart>
        </Box>

        {/* Only render the breakdown list when we have a real total */}
        {total > 0 &&
          pairs.map((pair, index) => {
            const pct = total > 0 ? (pair.value / total) * 100 : 0;
            return (
              <Stack key={index} direction="row" sx={{ alignItems: 'center', gap: 2, pb: 2 }}>
                <Stack sx={{ gap: 1, flexGrow: 1 }}>
                  <Stack direction="row" sx={{ justifyContent: 'space-between', alignItems: 'center', gap: 2 }}>
                    <Typography variant="body2" sx={{ fontWeight: '500' }}>
                      {pair.name}
                    </Typography>
                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                      {pct.toFixed(2)}%
                    </Typography>
                  </Stack>
                  <LinearProgress
                    variant="determinate"
                    aria-label="Number of users by country"
                    value={parseFloat(pct.toFixed(2))}
                    sx={{
                      [`& .${linearProgressClasses.bar}`]: {
                        backgroundColor: colors[index],
                      },
                    }}
                  />
                </Stack>
              </Stack>
            );
          })}
      </CardContent>
    </Card>
  );
}