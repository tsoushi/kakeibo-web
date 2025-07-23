import { useQuery } from 'urql'
import { graphql } from './gql/gql'

const GetUserDocument = graphql(/* GraphQL */`
  query GetUser {
    user {
      id
      name
    }
  }
`)

function App() {
  const [result] = useQuery({ query: GetUserDocument })

  const { data, fetching, error } = result
  if (fetching) return <p>Loading...</p>
  if (error) return <div>
    <p>GraphQLError: {error.graphQLErrors.map(e => e.message).join(', ')}</p>
  </div>

  return (
    <div>
      <h1>User Information</h1>
      {data?.user ? (
        <div>
          <p>ID: {data.user.id}</p>
          <p>Name: {data.user.name}</p>
        </div>
      ) : (
        <p>No user data found.</p>
      )}
    </div>
  )
}

export default App
