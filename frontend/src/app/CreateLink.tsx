import * as React from 'react';
import { alpha } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import getSignInSideTheme from './theme/getSignInSideTheme';
import TextField from '@mui/material/TextField';
import { createTheme, ThemeProvider, PaletteMode } from '@mui/material/styles';
import FormLabel from '@mui/material/FormLabel';
import FormControl from '@mui/material/FormControl';
import TextareaAutosize from '@mui/material/TextareaAutosize';
import SideMenu from './SideMenu';
import MainGrid from './MainGrid';
import axios from 'axios';
import Cookies from 'js-cookie';
import CustomCard from './CustomCard';
import { 
	Button,
	Typography,
	Divider
 } from '@mui/material';

export default function Dashboard(props: { disableCustomTheme?: boolean }) {
	const [mode, setMode] = React.useState<PaletteMode>('light');
	const [loggedIn, setLoggedIn] = React.useState(false);
	const [username, setUsername] = React.useState('');
	const [email, setEmail] = React.useState('');
	const [destinationError, setDestinationError] = React.useState(false);
	const [destinationErrorMessage, setDestinationErrorMessage] = React.useState('');
	const [titleError, setTitleError] = React.useState(false);
	const [titleErrorMessage, setTitleErrorMessage] = React.useState('');
	const [category, setCategory] = React.useState('option1');
	const [linkError, setLinkError] = React.useState(false);
	const [linkErrorMessage, setLinkErrorMessage] = React.useState('');
	const [ipError, setIpError] = React.useState(false);
	const [ipErrorMessage, setIpErrorMessage] = React.useState('');
	const SignInSideTheme = createTheme(getSignInSideTheme(mode));

	const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		if (destinationError || titleError || ipError) {
		  return;
		}
		console.log('submit');

		const data = new FormData(event.currentTarget);
		const referenced_link = data.get('destination');
		const title = data.get('title');
		const category = data.get('category');
		const allowed_ips = [];
		const black_listed_ips = [];
		if (category === 'option2') {
			const blacklisted = data.get('blacklisted');
			
			const ipToInt = (ip: string) => {
				return ip.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet, 10), 0) >>> 0;
			};

			const blacklistedIps = blacklisted.split('\n').map(ip => ip.trim()).filter(ip => ip !== '');
			blacklistedIps.forEach(ip => {
				black_listed_ips.push(ipToInt(ip));
			});
		}

		if (category === 'option3') {
			const whitelisted = data.get('whitelisted');

			const ipToInt = (ip: string) => {
				return ip.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet, 10), 0) >>> 0;
			};

			const whitelistedIps = whitelisted.split('\n').map(ip => ip.trim()).filter(ip => ip !== '');
			whitelistedIps.forEach(ip => {
				allowed_ips.push(ipToInt(ip));
			});	 
		}

		let access_mode = 0;
		if (category === 'option2') {
			access_mode = 2;
		} else if (category === 'option3') {
			access_mode = 1;
		}

		const authToken = Cookies.get('AuthToken');

		axios.post('http://localhost:3000/app/create-link', { 
		'referenced_link' : referenced_link, 
		'title' : title,
		'allowed_ips' : allowed_ips,
		'black_listed_ips' : black_listed_ips,
		'access_mode' : access_mode
		}, {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			redirectToLinks();
		}).catch(error => {
			setLinkError(true);
			setLinkErrorMessage('Failed to create link');
		});
	};


	const redirectToLinks = () => {
		window.location.href = '/app/links';
	}

	const validateInputs = () => {
		const destination = document.getElementById('destination') as HTMLInputElement;
		const title = document.getElementById('title') as HTMLInputElement;

		let isValid = true;

		const urlPattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
			'((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.?)+[a-z]{2,}|'+ // domain name
			'((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
			'(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
			'(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
			'(\\#[-a-z\\d_]*)?$','i'); // fragment locator

		if (!urlPattern.test(destination.value)) {
			setDestinationError(true);
			setDestinationErrorMessage('Invalid URL');
			isValid = false;
		} else {
			setDestinationError(false);
			setDestinationErrorMessage('');
		}

		if (title.value.length > 256) {
			setTitleError(true);
			setTitleErrorMessage('Title must be less than 256 characters');
			isValid = false;
		} else {
			setTitleError(false);
			setTitleErrorMessage('');
		}

		if (category === 'option2') {
			const blacklisted = document.getElementById('blacklisted') as HTMLInputElement;
			const ipPattern = new RegExp('^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$');
			
			const blacklistedIps = blacklisted.value.split('\n');
			for (const ip of blacklistedIps) {
				if (!ipPattern.test(ip.trim())) {
					setIpError(true);
					setIpErrorMessage('Invalid IP in blacklist');
					isValid = false;
					break;
				}
			}

			if (isValid) {
				setIpError(false);
				setIpErrorMessage('');
			}
		}

		if (category === 'option3') {
			const whitelisted = document.getElementById('whitelisted') as HTMLInputElement;
			const ipPattern = new RegExp('^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$');
			
			const whitelistedIps = whitelisted.value.split('\n');
			for (const ip of whitelistedIps) {
				if (!ipPattern.test(ip.trim())) {
					setIpError(true);
					setIpErrorMessage('Invalid IP in whitelist');
					isValid = false;
					break;
				}
			}

			if (isValid) {
				setIpError(false);
				setIpErrorMessage('');
			}
		}
		return isValid;
	};

	// This code only runs on the client side, to determine the system color preference
	React.useEffect(() => {
	  // Check if there is a preferred mode in localStorage
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
	  const systemPrefersDark = window.matchMedia(
		  '(prefers-color-scheme: dark)',
		).matches;
		setMode(systemPrefersDark ? 'dark' : 'light');
	  }, []);
  return (
    <ThemeProvider theme={SignInSideTheme}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: 'flex' }}>
        {/* Main content */}
		<SideMenu username={username} email={email} selectedItem='Links'/>
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
          <Stack
            spacing={2}
            sx={{
              alignItems: 'center',
              mx: 3,
              pb: 5,
              mt: { xs: 8, md: 1 },
            }}
          >
			<Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width:'100%', maxWidth:'600px'}}>
				<Typography variant="h3" component="div">
					Create Link
				</Typography>
			</Box>
				<Box
        component="form"
        onSubmit={handleSubmit}
        noValidate
        sx={{ display: 'flex', flexDirection: 'column', width: '100%', gap: 2, maxWidth: '600px' }}
      >
        <FormControl>
          <FormLabel htmlFor="destination">Destination</FormLabel>
          <TextField
            error={destinationError}
            helperText={destinationErrorMessage}
            id="destination"
            type="text"
            name="destination"
            placeholder="https://example.com"
            autoFocus
            required
            fullWidth
            variant="outlined"
            color={destinationError ? 'error' : 'primary'}
            sx={{ ariaLabel: 'destination' }}
          />
        </FormControl>
        <FormControl>
          <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
            <FormLabel htmlFor="title">Title (optional)</FormLabel>
          </Box>
          <TextField
            error={titleError}
            helperText={titleErrorMessage}
            name="title"
            placeholder=""
            type="text"
            id="title"
            autoFocus
            fullWidth
            variant="outlined"
            color={titleError ? 'error' : 'primary'}
          />
        </FormControl>
		<FormControl>
			<FormLabel htmlFor="category">Category</FormLabel>
			<TextField
				select
				id="category"
				name="category"
				variant="outlined"
				fullWidth
				SelectProps={{
					native: true,
				}}
				onChange={(event) => setCategory(event.target.value)}
			>
				<option value="option1">Accessible to anyone</option>
				<option value="option2">Blocked for blacklisted ips</option>
				<option value="option3">Only allowed to whitelisted ips</option>
			</TextField>
		</FormControl>
		{category === 'option2' && (
			<>
			<FormControl>
				<FormLabel htmlFor="blacklisted">BlackListed Ips</FormLabel>
					<TextareaAutosize
						id="blacklisted"
						name="blacklisted"
						minRows={5}
					/>
			</FormControl>
			<Typography color="error" sx={{ marginTop: 'none' }}>
				{ipErrorMessage}
			</Typography>
			</>
		)}

		{category === 'option3' && (
			<>
			<FormControl>
				<FormLabel htmlFor="whitelisted">WhiteListed Ips</FormLabel>
					<TextareaAutosize
						id="whitelisted"
						name="whitelisted"
						minRows={5}
					/>
			</FormControl>
			<Typography color="error" sx={{ marginTop: 'none' }}>
				{ipErrorMessage}
			</Typography>
			</>
		)}
		<Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
			<Button variant="outlined" color="secondary" onClick={redirectToLinks}>
				Cancel
			</Button>
			<Button variant="contained" color="primary" type="submit" onClick={validateInputs}>
				Submit
			</Button>
		</Box>
		{linkError && (
			<Typography color="error" sx={{ textAlign: 'center' }}>
				{linkErrorMessage}
			</Typography>
		)}
      </Box>
	</Stack>
	</Box>
	</Box>
    </ThemeProvider>
  );
}