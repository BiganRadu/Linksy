import * as React from 'react';
import { styled } from '@mui/material/styles';
import Avatar from '@mui/material/Avatar';
import MuiDrawer, { drawerClasses } from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import MenuContent from './MenuContent';
import IconButton from '@mui/material/IconButton';
import LogoutIcon from '@mui/icons-material/Logout';
import axios from 'axios';
import Cookies from 'js-cookie';

const drawerWidth = 240;

const Drawer = styled(MuiDrawer)({
  width: drawerWidth,
  flexShrink: 0,
  boxSizing: 'border-box',
  mt: 10,
  [`& .${drawerClasses.paper}`]: {
    width: drawerWidth,
    boxSizing: 'border-box',
  },
});

interface SideMenuInfo {
	username: string;
	email: string;
	selectedItem: string;
}

const logout = async () => {
	try {
		const response = await axios.get('http://localhost:3000/member/logout', {
			headers: {
				'AuthToken': Cookies.get('AuthToken')
			}
		});

		if (response.status === 200) {
			Cookies.remove('AuthToken');
			window.location.href = '/sign-in';
		}
	} catch (error) {
		console.error('Error logging out:', error);
	}
}

export default function SideMenu({ username, email, selectedItem} : SideMenuInfo) {
  return (
    <Drawer
      variant="permanent"
      sx={{
        display: { xs: 'none', md: 'block' },
        [`& .${drawerClasses.paper}`]: {
          backgroundColor: 'background.paper',
        },
      }}
    >
      <MenuContent selectedItem= {selectedItem}/>
      <Stack
        direction="row"
        sx={{
          p: 2,
          gap: 1,
          alignItems: 'center',
          borderTop: '1px solid',
          borderColor: 'divider',
        }}
      >
        <Avatar
          sizes="small"
          alt={username}
          src="/static/images/avatar/8.jpg"
          sx={{ width: 36, height: 36 }}
        />
        <Box sx={{ mr: 'auto' }}>
          <Typography variant="body2" sx={{ fontWeight: 500, lineHeight: '16px' }}>
            {username}
          </Typography>
          <Typography variant="caption" sx={{ color: 'text.secondary' }}>
            {email}
          </Typography>
        </Box>

		<IconButton size="small" edge="end" color="inherit" onClick={logout}>
			<LogoutIcon />
		</IconButton>
      </Stack>
    </Drawer>
  );
}
