import * as React from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Stack from '@mui/material/Stack';
import HomeRoundedIcon from '@mui/icons-material/HomeRounded';
import AnalyticsRoundedIcon from '@mui/icons-material/AnalyticsRounded';
import LinkIcon from '@mui/icons-material/Link';
import QrCodeIcon from '@mui/icons-material/QrCode';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';
import { redirect } from 'react-router-dom';

const mainListItems = [
  { text: 'Home', icon: <HomeRoundedIcon />, redirect: '/app'},
  { text: 'Links', icon: <LinkIcon />, redirect: '/app/links' },
  { text: 'QR Codes', icon: <QrCodeIcon />, redirect: '/app/qrcodes' },
  { text: 'Analytics', icon: <AnalyticsRoundedIcon />, redirect: '/app/analytics' },
];

function handleRedirect(path: string) {
	if (window.location.pathname !== path) {
		window.location.href = path
	}
}

const secondaryListItems = [
  { text: 'Settings', icon: <SettingsRoundedIcon />, redirect: '/app/settings' },
];

interface MenuContentProps {
	selectedItem: string;
}

function getSelectedIndex(items: { text: string; icon: JSX.Element }[], selectedItem: string): number {
	return items.findIndex(item => item.text.toLowerCase() === selectedItem.toLowerCase());
}
export default function MenuContent({ selectedItem }: MenuContentProps) {
	const selectedIndex = getSelectedIndex([...mainListItems, ...secondaryListItems], selectedItem);
	
	return (
		<Stack sx={{ flexGrow: 1, p: 1, justifyContent: 'space-between' }}>
			<List dense>
				{mainListItems.map((item, index) => (
					<ListItem key={index} disablePadding sx={{ display: 'block' }}>
						<ListItemButton selected={index === selectedIndex} onClick={() => handleRedirect(item.redirect)}>
							<ListItemIcon>{item.icon}</ListItemIcon>
							<ListItemText primary={item.text} />
						</ListItemButton>
					</ListItem>
				))}
			</List>

			<List dense>
				{secondaryListItems.map((item, index) => (
					<ListItem key={index + mainListItems.length} disablePadding sx={{ display: 'block' }}>
						<ListItemButton selected={index + mainListItems.length === selectedIndex} onClick={() => handleRedirect(item.redirect)}>
							<ListItemIcon>{item.icon}</ListItemIcon>
							<ListItemText primary={item.text} />
						</ListItemButton>
					</ListItem>
				))}
			</List>
		</Stack>
	);
}
