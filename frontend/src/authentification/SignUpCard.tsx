import * as React from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import MuiCard from '@mui/material/Card';
import FormLabel from '@mui/material/FormLabel';
import FormControl from '@mui/material/FormControl';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import axios from 'axios';

import { styled } from '@mui/material/styles';



const Card = styled(MuiCard)(({ theme }) => ({
	display: 'flex',
	flexDirection: 'column',
	alignSelf: 'center',
	width: '100%',
	padding: theme.spacing(4),
	gap: theme.spacing(2),
	boxShadow:
	  'hsla(220, 30%, 5%, 0.05) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.05) 0px 15px 35px -5px',
	[theme.breakpoints.up('sm')]: {
	  width: '450px',
	},
	...theme.applyStyles('dark', {
	  boxShadow:
		'hsla(220, 30%, 5%, 0.5) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.08) 0px 15px 35px -5px',
	}),
  }));
  
export default function SignUpCard() {
	const [nameError, setNameError] = React.useState(false);
	const [nameErrorMessage, setNameErrorMessage] = React.useState('');
	const [emailError, setEmailError] = React.useState(false);
	const [emailErrorMessage, setEmailErrorMessage] = React.useState('');
	const [passwordError, setPasswordError] = React.useState(false);
	const [passwordErrorMessage, setPasswordErrorMessage] = React.useState('');
	const [signUpError, setSignUpError] = React.useState(false);
	const [signUpErrorMessage, setSignUpErrorMessage] = React.useState('');
  
	const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
	  event.preventDefault();
	  if (emailError || passwordError || nameError) {
		return;
	  }
	  const data = new FormData(event.currentTarget);
	  const email = data.get('email');
	  const password = data.get('password');
	  const username = data.get('name');

	axios.post('http://localhost:3000/member/register', {
		email,
		password,
		username,
	})
	.then(response => {
		window.location.href = '/sign-in';
	})
	.catch(error => {
		console.log(error.response.data['error']);
		setSignUpError(true);
		setSignUpErrorMessage(error.response.data['error']);
	});
	};
  
	const validateInputs = () => {
	  const email = document.getElementById('email') as HTMLInputElement;
	  const password = document.getElementById('password') as HTMLInputElement;
	  const name = document.getElementById('name') as HTMLInputElement;
  
	  let isValid = true;
  
	  if (!email.value || !/\S+@\S+\.\S+/.test(email.value)) {
		setEmailError(true);
		setEmailErrorMessage('Please enter a valid email address.');
		isValid = false;
	  } else {
		setEmailError(false);
		setEmailErrorMessage('');
	  }
  
	  if (!password.value || password.value.length < 6) {
		setPasswordError(true);
		setPasswordErrorMessage('Password must be at least 6 characters long.');
		isValid = false;
	  } else {
		setPasswordError(false);
		setPasswordErrorMessage('');
	  }

	  if (!name.value || name.value.length < 1) {
		setNameError(true);
		setNameErrorMessage('Name is required.');
		isValid = false;
	  } else {
		setNameError(false);
		setNameErrorMessage('');
	  }
  
	  return isValid;
	};
  
	return (
	  <Card variant="outlined">
		<Box sx={{ display: { xs: 'flex', md: 'none' } }}>
		</Box>
		<Typography
		  component="h1"
		  variant="h4"
		  sx={{ width: '100%', fontSize: 'clamp(2rem, 10vw, 2.15rem)' }}
		>
		  Sign up
		</Typography>
		<Box
		  component="form"
		  onSubmit={handleSubmit}
		  noValidate
		  sx={{ display: 'flex', flexDirection: 'column', width: '100%', gap: 2 }}
		>
		<FormControl>
			<FormLabel htmlFor="email">Username</FormLabel>
			<TextField
			  error={nameError}
			  helperText={nameErrorMessage}
			  id="name"
			  type="name"
			  name="name"
			  placeholder="Jon Snow"
			  autoComplete="name"
			  autoFocus
			  required
			  fullWidth
			  variant="outlined"
			  color={nameError ? 'error' : 'primary'}
			  sx={{ ariaLabel: 'name' }}
			/>
		</FormControl>
		  <FormControl>
			<FormLabel htmlFor="email">Email</FormLabel>
			<TextField
			  error={emailError}
			  helperText={emailErrorMessage}
			  id="email"
			  type="email"
			  name="email"
			  placeholder="your@email.com"
			  autoComplete="email"
			  autoFocus
			  required
			  fullWidth
			  variant="outlined"
			  color={emailError ? 'error' : 'primary'}
			  sx={{ ariaLabel: 'email' }}
			/>
		  </FormControl>
		  <FormControl>
			<Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
			  <FormLabel htmlFor="password">Password</FormLabel>
			</Box>
			<TextField
			  error={passwordError}
			  helperText={passwordErrorMessage}
			  name="password"
			  placeholder="••••••"
			  type="password"
			  id="password"
			  autoComplete="current-password"
			  autoFocus
			  required
			  fullWidth
			  variant="outlined"
			  color={passwordError ? 'error' : 'primary'}
			/>
		  </FormControl>
		  <Button type="submit" fullWidth variant="contained" onClick={validateInputs}>
			Sign up
		  </Button>
		{signUpError && (
			<Typography color="error" variant="body2" sx= {{textAlign: 'center'}}>
				{signUpErrorMessage}
			</Typography>
		)}
		  <Typography sx={{ textAlign: 'center' }}>
			Already have an account?{' '}
			<span>
			  <Link
				href="/sign-in"
				variant="body2"
				sx={{ alignSelf: 'center' }}
			  >
				Sign in
			  </Link>
			</span>
		  </Typography>
		</Box>
	  </Card>
	);
  }
