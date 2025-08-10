import * as React from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { GridCellParams, GridColDef } from '@mui/x-data-grid';
import { SparkLineChart } from '@mui/x-charts/SparkLineChart';
import axios from 'axios';
import Cookies from 'js-cookie';

type SparkLineData = number[];

function getNameOfTheDays(startTimestamp: number, endTimestamp: number) : string[] {
	const days = [];
	const today = new Date(endTimestamp * 1000);
	let numOfDays = Math.floor((endTimestamp - startTimestamp) / (60 * 60 * 24));
	for (let i = numOfDays; i >= 1; i--) {
		const d = new Date(today);
		d.setDate(today.getDate() - i);
		const monthName = d.toLocaleDateString('en-US', { month: 'short' });
		days.push(`${monthName} ${d.getDate()}`);
	}
	return days;
}

function renderSparklineCell(params: GridCellParams<SparkLineData, any>) {
  const { value, colDef } = params;

  if (!value || value.length === 0) {
    return null;
  }

  return (
    <div style={{ display: 'flex', alignItems: 'center', height: '100%' }}>
      <SparkLineChart
        data={value}
        width={colDef.computedWidth || 100}
        height={32}
        plotType="bar"
        showHighlight
        showTooltip
        color="hsl(210, 98%, 42%)"
        xAxis={{
          scaleType: 'band',
          data: params.row.days, // Use the days from the row data
        }}
      />
    </div>
  );
}

const columns: GridColDef[] = [
  { field: 'title', headerName: 'Link Title', flex: 1.5, minWidth: 200 },
  {
    field: 'total',
    headerName: 'Accesses',
    headerAlign: 'right',
    align: 'right',
    flex: 1,
    minWidth: 80,
  },
  {
    field: 'country',
    headerName: 'Best Country',
    headerAlign: 'right',
    align: 'right',
    flex: 1,
    minWidth: 100,
  },
  {
    field: 'platform',
    headerName: 'Best Platform',
    headerAlign: 'right',
    flex: 1,
    minWidth: 100,
  },
  {
    field: 'sessions',
    headerName: 'Daily Accesses',
    flex: 1,
    minWidth: 150,
    renderCell: renderSparklineCell,
  },
];


export default function CustomizedDataGrid({startTimestamp, endTimestamp}: {startTimestamp: number, endTimestamp: number}) {
	const [rows, setRows] = React.useState([]);
	const days = getNameOfTheDays(startTimestamp, endTimestamp);

	// Fetch data from the server
	React.useEffect(() => {
		const fetchData = async () => {
			const authToken = Cookies.get('AuthToken');
				axios.get(`http://localhost:3000/app/analytics?chart_code=links&start=${startTimestamp}&end=${endTimestamp}`, {
				headers: {
					AuthToken: authToken,
				},
		}).then(response => {
			setRows(response.data.links);
		}).catch(error => {
			console.log("Error fetching analytics data:", error);
		});
	};

		fetchData();
	}, [startTimestamp, endTimestamp]);

	// Add the days to each row
	for (let i = 0; i < rows.length; i++) {
		rows[i].days = days;
	}
  return (
    <DataGrid
      rows={rows}
      columns={columns}
      getRowClassName={(params) =>
        params.indexRelativeToCurrentPage % 2 === 0 ? 'even' : 'odd'
      }
      initialState={{
        pagination: { paginationModel: { pageSize: 20 } },
      }}
      pageSizeOptions={[10, 20, 50]}
      disableColumnResize
      density="compact"
      slotProps={{
        filterPanel: {
          filterFormProps: {
            logicOperatorInputProps: {
              variant: 'outlined',
              size: 'small',
            },
            columnInputProps: {
              variant: 'outlined',
              size: 'small',
              sx: { mt: 'auto' },
            },
            operatorInputProps: {
              variant: 'outlined',
              size: 'small',
              sx: { mt: 'auto' },
            },
            valueInputProps: {
              InputComponentProps: {
                variant: 'outlined',
                size: 'small',
              },
            },
          },
        },
      }}
    />
  );
}