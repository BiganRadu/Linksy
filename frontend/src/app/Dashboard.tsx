import * as React from 'react';
import { alpha } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import getSignInSideTheme from '../theme/getSignInSideTheme';
import { createTheme, ThemeProvider, PaletteMode } from '@mui/material/styles';
import SideMenu from './SideMenu';
import MainGrid from './MainGrid';
import axios from 'axios';
import Cookies from 'js-cookie';
import CircularProgress from '@mui/material/CircularProgress'; // added

export default function Dashb(props: { disableCustomTheme?: boolean }) {
  const [mode, setMode] = React.useState<PaletteMode>('light');
  const [username, setUsername] = React.useState('');
  const [email, setEmail] = React.useState('');
  const [loading, setLoading] = React.useState(true); // added
  const SignInSideTheme = createTheme(getSignInSideTheme(mode));

  React.useEffect(() => {
    const fetchMemberInfo = async () => {
      const authToken = Cookies.get('AuthToken');
      setLoading(true); // added
      axios
        .get('https://linksy-mhe5.onrender.com/app/member-info', {
          headers: {
            AuthToken: authToken,
          },
        })
        .then((response) => {
          setUsername(response.data.username);
          setEmail(response.data.email);
        })
        .catch((error) => {
          window.location.href = '/sign-in';
        })
        .finally(() => {
          setLoading(false); // added
        });
    };

    fetchMemberInfo();
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setMode(systemPrefersDark ? 'dark' : 'light');
  }, []);

  return (
    <ThemeProvider theme={SignInSideTheme}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: 'flex', minHeight: '100vh' }}>
        {loading ? (
          <Box
            sx={{
              display: 'flex',
              flexGrow: 1,
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <CircularProgress />
          </Box>
        ) : (
          <>
            <SideMenu username={username} email={email} selectedItem="Home" />
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
                  mt: { xs: 8, md: 0 },
                }}
              >
                <MainGrid />
              </Stack>
            </Box>
          </>
        )}
      </Box>
    </ThemeProvider>
  );
}
