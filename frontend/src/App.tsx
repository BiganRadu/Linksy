import {
	BrowserRouter as Router,
	Routes,
	Route,
	Link
} from 'react-router-dom'
import Redirecter from './Redirecter'
import SignInSide from './authentification/SignInSide'
import SignUpSide from './authentification/SignUpSide'
import Dashboard from './app/Dashboard'
import LinksPage from './app/LinksPage'
import CreateLink from './app/CreateLink'
import EditLink from './app/EditLink'
function App() {

return (
	<Router>
		<Routes>
			<Route path="/app" element={<Dashboard/>} />
			<Route path="/app/links" element={<LinksPage/>} />
			<Route path="/app/create-link" element={<CreateLink/>} />
			<Route path="/app/edit-link/:linkID" element={<EditLink />} />
			<Route path="/sign-up" element={<SignUpSide />}/>
			<Route path="/sign-in" element= {<SignInSide/>}/>
			<Route path="/:linkID" element={<Redirecter />} />
		</Routes>
	</Router>
)
}

export default App
