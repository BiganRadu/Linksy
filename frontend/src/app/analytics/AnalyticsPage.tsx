import * as React from 'react';
import { alpha } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Grid from '@mui/material/Grid';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import getSignInSideTheme from '../../theme/getSignInSideTheme';
import { createTheme, ThemeProvider, PaletteMode } from '@mui/material/styles';
import SideMenu from '../SideMenu';
import axios from 'axios';
import Cookies from 'js-cookie';
import { Typography } from '@mui/material';
import SessionsChart from './SessionsChart';
import RoundChart from './RoundChart';
import CustomizedDataGrid from './CustomizedDataGrid';
import dayjs from 'dayjs';

export default function Dashboard(props: { disableCustomTheme?: boolean }) {
	const [mode, setMode] = React.useState<PaletteMode>('light');
	const [loggedIn, setLoggedIn] = React.useState(false);
	const [username, setUsername] = React.useState('');
	const [email, setEmail] = React.useState('');
	const SignInSideTheme = createTheme(getSignInSideTheme(mode));
	let toDate = new Date();
	let fromDate = new Date();
	fromDate.setDate(toDate.getDate() - 29);
	const [startTimestamp, setStartTimestamp] = React.useState(Math.floor(fromDate.getTime() / 1000));
	const [endTimestamp, setEndTimestamp] = React.useState(Math.floor(toDate.getTime() / 1000))
	React.useEffect(() => {

	const fetchMemberInfo = async () => {

		const authToken = Cookies.get('AuthToken');
		axios.get('http://localhost:3000/app/member-info', {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			setUsername(response.data.username);
			setEmail(response.data.email);
			setLoggedIn(true);
		}).catch(error => {
			window.location.href = '/sign-in';
		});
	};

	fetchMemberInfo();

	  // Check if the user prefers dark mode
	  const systemPrefersDark = window.matchMedia(
		  '(prefers-color-scheme: dark)',
		).matches;
		setMode(systemPrefersDark ? 'dark' : 'light');
	  }, []);
  return (
	<ThemeProvider theme={SignInSideTheme}>
	  <CssBaseline enableColorScheme />
	  <Box sx={{ display: 'flex' }}>
		<SideMenu username={username} email={email} selectedItem='Analytics'/>
		<Box
		  component="main"
		  sx={(theme) => ({
			flexGrow: 1,
			backgroundColor: theme.vars
			  ? `rgba(${theme.vars.palette.background.defaultChannel} / 1)`
			  : alpha(theme.palette.background.default, 1),
			overflow: 'auto',
		  })}
		>
		<Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 2, ml: 3, mt: 2, mr: 3 }}>
			<Typography component="h2" variant="h3">
				Analytics Page
			</Typography>
			<Stack direction="row" spacing={2}>
				<LocalizationProvider dateAdapter={AdapterDayjs}>
					<DatePicker
						label="From"
						value={dayjs.unix(startTimestamp)}
						onChange={(e) => {
							setStartTimestamp(Math.floor(e?.toDate().getTime() / 1000));
						}}
					/>
				</LocalizationProvider>
				<LocalizationProvider dateAdapter={AdapterDayjs}>
					<DatePicker
						label="To"
						value={dayjs.unix(endTimestamp)}
						onChange={(e) => {
							setEndTimestamp(Math.floor(e?.toDate().getTime() / 1000));
						}}
					/>
				</LocalizationProvider>
			</Stack>
		</Stack>

		{ /* Main content of the analytics page */}
		<Box sx= {{ml:3, mr:3}}>
			<Grid size={{ xs: 12, md: 6 }}>
          		<SessionsChart startTimestamp={startTimestamp} endTimestamp={endTimestamp} />
        	</Grid>
        <Grid size={{ xs: 12, lg: 3 }} sx= {{mt: 2, mb: 2}}>
          <Stack gap={2} direction={{ xs: 'row', sm: 'row', lg: 'row' }}>
			<Grid size={{ xs: 12, sm: 6, md: 6, lg: 12 }} sx={{width: '50%'}}>
            	<RoundChart chart_code='countries' startTimestamp={startTimestamp} endTimestamp={endTimestamp} />
			</Grid>
			<Grid size={{ xs: 12, sm: 6, md: 6, lg: 12 }} sx={{width: '50%'}}>
				<RoundChart chart_code='platforms' startTimestamp={startTimestamp} endTimestamp={endTimestamp} />
			</Grid>
          </Stack>
		</Grid>
		<Grid size={{ xs: 12, lg: 9 }}>
          <CustomizedDataGrid startTimestamp={startTimestamp} endTimestamp={endTimestamp} />
        </Grid>
		</Box>
		</Box>
	  </Box>
	</ThemeProvider>
  );
}
