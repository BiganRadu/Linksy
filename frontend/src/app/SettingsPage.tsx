import * as React from 'react';
import { alpha } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import getSignInSideTheme from '../theme/getSignInSideTheme';
import { createTheme, ThemeProvider, PaletteMode } from '@mui/material/styles';
import SideMenu from './SideMenu';
import axios from 'axios';
import Cookies from 'js-cookie';
import {
	TextField,
	FormControl,
 } from '@mui/material';
import {
	Button,
	Typography,
 } from '@mui/material';

export default function Dashboard(props: { disableCustomTheme?: boolean }) {
	const [mode, setMode] = React.useState<PaletteMode>('light');
	const [username, setUsername] = React.useState('');
	const [tempName, setTempName] = React.useState('');
	const [email, setEmail] = React.useState('');
	const [passwordError, setPasswordError] = React.useState('');
	const [alertMessage, setAlertMessage] = React.useState('');
	const SignInSideTheme = createTheme(getSignInSideTheme(mode));

	const validateInputs = (oldPassword: HTMLInputElement, newPassword: HTMLInputElement, confirmNewPassword: HTMLInputElement) => {

		if (!oldPassword || !newPassword || !confirmNewPassword) {
			setPasswordError('All fields are required.');
			return false;
		}
		if (newPassword.value !== confirmNewPassword.value) {
			setPasswordError('New passwords do not match.');
			return false;
		}
		if (newPassword.value.length < 6) {
			setPasswordError('New password must be at least 6 characters long.');
			return false;
		}
		return true;
	};

	const handleChangeName = () => {
		const nameInput = document.getElementById('name') as HTMLInputElement;
		if (!nameInput || !nameInput.value) {
			setAlertMessage('Name cannot be empty.');
			return;
		}
		const authToken = Cookies.get('AuthToken');
		axios.post('http://localhost:3000/member/change-name', {
			"new_name": nameInput.value,
		}, {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			setAlertMessage('Name changed successfully.');
			setUsername(nameInput.value);
			nameInput.value = '';
		}).catch(error => {
			console.log("Error changing name:", error);
		});
	}

	const handleChangePassword = () => {
		const oldPassword = document.getElementById('current_password') as HTMLInputElement;
		const newPassword = document.getElementById('new_password') as HTMLInputElement;
		const confirmNewPassword = document.getElementById('confirm_new_password') as HTMLInputElement;
		if (!validateInputs(oldPassword, newPassword, confirmNewPassword)) {
			return;
		}
		const authToken = Cookies.get('AuthToken');
		axios.post('http://localhost:3000/member/change-password', {
			"old_password": oldPassword.value,
			"new_password": newPassword.value,
		}, {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			setAlertMessage('Password changed successfully.');
			setPasswordError('');
			
			oldPassword.value = '';
			newPassword.value = '';
			confirmNewPassword.value = '';
		}).catch(error => {
			setPasswordError(error.response.data.error || 'An error occurred while changing the password.');
		});
	};

	React.useEffect(() => {

	const fetchMemberInfo = async () => {
		const authToken = Cookies.get('AuthToken');
		axios.get('http://localhost:3000/app/member-info', {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			setUsername(response.data.username);
			setTempName(response.data.username);
			setEmail(response.data.email);
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
		<SideMenu username={username} email={email} selectedItem='Settings'/>
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
			<Box sx={{ml:3}}>
		    <Typography component="h2" variant="h3" sx={{ mb: 2, mt: 2 }}>
			  Profile Settings
		    </Typography>
			<Typography variant="h5" >
				Display Name
			</Typography>
			<FormControl sx= {{width: '50%'}}>
			<TextField
				id="name"
				type="text"
				name="name"
				autoComplete="name"
				autoFocus
				required
				fullWidth
				variant="outlined"
				value={tempName}
				onChange={(e) => setTempName(e.target.value)}
				sx={{ ariaLabel: 'name', width:'100%' }}
			/>
			</FormControl>
			<Box sx={{ display: 'flex', justifyContent: 'flex-start', width: '50%' }}>
			<Button
				variant="contained"
				color="primary"
				sx={{
					mt: 2,
					cursor: tempName === username ? 'not-allowed' : 'pointer',
				}}
				onClick={() => {
					if (tempName === username) {
						return;
					}
					handleChangeName();
				}}
			>
				Change Name
			</Button>
			</Box>
			<Typography variant="h5" sx={{ mt: 2 }}>
				Email
			</Typography>
			<Typography id="top" variant="body1">
				{email}
			</Typography>
			<Typography variant="h4">
				Change Password
			</Typography>
			<Typography variant="body1" sx={{ mb: 2 }}>
				You will be required to login after changing your password
			</Typography>
			<Typography variant="h5">
				Current Password
			</Typography>
			<FormControl sx= {{width: '50%'}}>
			<TextField
				id="current_password"
				type="password"
				name="current_password"
				autoComplete="current_password"
				required
				fullWidth
				variant="outlined"
				sx={{ ariaLabel: 'current_password', width:'100%' }}
			/>
			</FormControl>
			<Typography variant="h5" sx={{ mt: 2 }}>
				New Password
			</Typography>
			<FormControl sx= {{width: '50%'}}>
			<TextField
				id="new_password"
				type="password"
				name="new_password"
				autoComplete="new_password"
				required
				fullWidth
				variant="outlined"
				sx={{ ariaLabel: 'new_password', width:'100%' }}
			/>
			</FormControl>
			<Typography variant="h5" sx={{ mt: 2 }}>
				Confirm New Password
			</Typography>
			<FormControl sx= {{width: '50%'}}>
			<TextField
				id="confirm_new_password"
				type="password"
				name="confirm_new_password"
				autoComplete="confirm_new_password"
				required
				fullWidth
				variant="outlined"
				sx={{ ariaLabel: 'confirm_new_password', width:'100%' }}
			/>
			</FormControl>
			{passwordError && (
				<Typography variant="body2" sx={{ color: 'error.main', mt: 1 }}>
					{passwordError}
				</Typography>
			)}
			<Box sx={{ display: 'flex', justifyContent: 'flex-start', width: '50%' }}>
				{/* The Change Password button will be left-aligned by default inside this Box */}
			<Button
				variant="contained"
				color="primary"
				sx={{ mt: 2}}
				onClick={() => {
					handleChangePassword();
				}}
				>
				Change Password
			</Button>
			
			{alertMessage && (
				<Box
					sx={{
						position: 'fixed',
						top: 0,
						left: 0,
						width: '100vw',
						height: '100vh',
						display: 'flex',
						alignItems: 'center',
						justifyContent: 'center',
						zIndex: 1300,
						backgroundColor: 'rgba(0,0,0,0.3)',
					}}
				>
					<Box
						sx={{
							bgcolor: 'background.paper',
							p: 4,
							borderRadius: 2,
							boxShadow: 6,
							minWidth: 300,
							textAlign: 'center',
						}}
					>
						<Typography variant="h6" sx={{ mb: 2 }}>
							{alertMessage}
						</Typography>
						
						<Button
							variant="contained"
							onClick={() => setAlertMessage('')}
						>
							OK
						</Button>
					</Box>
				</Box>
			)}
			</Box>

			</Box>
		</Box>
	  </Box>
	</ThemeProvider>
  );
}