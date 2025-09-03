import * as React from 'react';
import { alpha } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import getSignInSideTheme from '../theme/getSignInSideTheme';
import { createTheme, ThemeProvider, PaletteMode } from '@mui/material/styles';
import SideMenu from './SideMenu';
import axios from 'axios';
import Cookies from 'js-cookie';
import CustomCard from './CustomCard';
import FacebookIcon from '@mui/icons-material/Facebook';
import TwitterIcon from '@mui/icons-material/Twitter';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import CircularProgress from '@mui/material/CircularProgress';
import {
	Button,
	IconButton,
	Typography,
	Divider
 } from '@mui/material';

export default function Dashboard(props: { disableCustomTheme?: boolean }) {
	const [mode, setMode] = React.useState<PaletteMode>('light');
	const [loggedIn, setLoggedIn] = React.useState(false);
	const [username, setUsername] = React.useState('');
	const [email, setEmail] = React.useState('');
	const [links, setLinks] = React.useState([]);
	const [linksError, setLinksError] = React.useState('');
	const [confirmation, setConfirmation] = React.useState(false);
	const [corfirmationId, setConfirmationId] = React.useState('');
	const [shareDialog, setShareDialog] = React.useState(false);
	const [loading, setLoading] = React.useState(true);
	const SignInSideTheme = createTheme(getSignInSideTheme(mode));

	const redirectToCreateLink = () => {
		window.location.href = '/app/create-link';
	}

	const redirectToEditLink = (id: string) => {
		window.location.href = `/app/edit-link/${id}`;
	}

	const handleDelete = (id: string) => {
		setConfirmation(true);
		setConfirmationId(id);
	}

	const handleShare = () => {
		setShareDialog(true);
	}

	const submitDelete = () => {
		const authToken = Cookies.get('AuthToken');
		axios.delete(`https://linksy-mhe5.onrender.com/app/delete-link?link_id=${corfirmationId}`, {
			headers: {
				AuthToken: authToken,
			},
		}).then(response => {
			console.log(response);
			window.location.reload();
		}).catch(error => {
			console.log(error);
		})
	}

	React.useEffect(() => {
	  const authToken = Cookies.get('AuthToken');
	  setLoading(true);

	  const fetchMemberInfo = () =>
		axios
		  .get('https://linksy-mhe5.onrender.com/app/member-info', {
			headers: { AuthToken: authToken },
		  })
		  .then((response) => {
			setUsername(response.data.username);
			setEmail(response.data.email);
			setLoggedIn(true);
		  })
		  .catch(() => {
			window.location.href = '/sign-in';
		  });

	  const fetchMemberLinks = () =>
		axios
		  .get('https://linksy-mhe5.onrender.com/app/member-links', {
			headers: { AuthToken: authToken },
		  })
		  .then((response) => {
			setLinks(response.data.links);
		  })
		  .catch(() => {
			setLinksError('Could not fetch links');
		  });

	  Promise.allSettled([fetchMemberInfo(), fetchMemberLinks()])
		.finally(() => setLoading(false));

	  const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
	  setMode(systemPrefersDark ? 'dark' : 'light');
	}, []);

	React.useEffect(() => {
		if (confirmation || shareDialog) {
			document.body.style.overflow = 'hidden'; // Disable scrolling
		} else {
			document.body.style.overflow = 'auto'; // Enable scrolling
		}
	
		// Cleanup function to restore scrolling when component unmounts or state changes
		return () => {
			document.body.style.overflow = 'auto';
		};
	}, [confirmation, shareDialog]);

	return (
		<ThemeProvider theme={SignInSideTheme}>
			<CssBaseline enableColorScheme />
			<Box sx={{ display: 'flex', minHeight: '100vh' }}>
				{loading ? (
					<Box sx={{ display: 'flex', flexGrow: 1, alignItems: 'center', justifyContent: 'center' }}>
						<CircularProgress />
					</Box>
				) : (
					<>
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
									Links Page
								</Typography>
								<Button variant="contained" size="small" onClick={redirectToCreateLink}>
									Create Link
								</Button>
							</Box>
							<Divider />
							{links === undefined || links === null ? (
								<p>No links available</p>
							) : (
								links.map((link, index) => (
									<CustomCard
										key={index}
										title={link.title}
										icon_url={link.icon}
										link_id={link.id}
										original_url={link.referenced_link}
										created_at={link.created_at}
										onEdit={() => redirectToEditLink(link.id)}
										onDelete={() => handleDelete(link.id)}
										onShare={handleShare}
									/>
								))
							)}
							</Stack>
						</Box>
					</>
				)}
			</Box>
	{confirmation && (
		<Box
			sx={{
				position: 'fixed',
				top: 0,
				left: 0,
				width: '100%',
				height: '100%',
				backgroundColor: 'rgba(0, 0, 0, 0.5)',
				display: 'flex',
				justifyContent: 'center',
				alignItems: 'center',
				zIndex: 1300,
			}}
		>
			<Box
				sx={{
					backgroundColor: 'blue',
					padding: 4,
					borderRadius: 2,
					boxShadow: 24,
					textAlign: 'center',
				}}
			>
				<Typography variant="h6" component="div" gutterBottom>
					Are you sure you want to delete this link?
				</Typography>
				<Button
					variant="contained"
					color="primary"
					onClick={() => {
						submitDelete();
						setConfirmation(false);
						setConfirmationId('');
					}}
					sx={{ mr: 2 }}
				>
					Yes
				</Button>
				<Button
					variant="outlined"
					color="secondary"
					onClick={() => {
						setConfirmation(false)
						setConfirmationId('');
					}}
				>
					No
				</Button>
			</Box>
		</Box>
	)}
	{shareDialog && (
		<Box
			sx={{
				position: 'fixed',
				top: 0,
				left: 0,
				width: '100%',
				height: '100%',
				backgroundColor: 'rgba(0, 0, 0, 0.5)',
				display: 'flex',
				justifyContent: 'center',
				alignItems: 'center',
				zIndex: 1300,
			}}
		>
			<Box
				sx={{
					backgroundColor: 'blue',
					padding: 4,
					borderRadius: 2,
					boxShadow: 24,
					position: 'relative',
					textAlign: 'center',
				}}
			>
				<Button
					sx={{
						position: 'absolute',
						top: 8,
						right: 8,
						backgroundColor: 'transparent',
						'&:hover': {
							backgroundColor: 'transparent',
						},
						'&:focus': {
							outline: 'none',
						},
					}}
					onClick={() => setShareDialog(false)}
				>
					X
				</Button>
				<Typography variant="h6" component="div" gutterBottom>
					Share this link
				</Typography>
				<Box sx={{ display: 'flex', justifyContent: 'center', gap: 2 }}>
					<IconButton
						color="primary"
						href="https://www.facebook.com/sharer/sharer.php?u=YOUR_LINK_HERE"
						target="_blank"
					>
						<FacebookIcon />
					</IconButton>
					<IconButton
						color="primary"
						href="https://twitter.com/intent/tweet?url=YOUR_LINK_HERE"
						target="_blank"
					>
						<TwitterIcon />
					</IconButton>
					<IconButton
						color="primary"
						href="https://www.linkedin.com/shareArticle?mini=true&url=YOUR_LINK_HERE"
						target="_blank"
					>
						<LinkedInIcon />
					</IconButton>
				</Box>
			</Box>
		</Box>
	)}
    </ThemeProvider>
  );
}