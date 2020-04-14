import React from 'react';
import { Box, useSnackbar } from 'pouncejs';
import SubmitButton from 'Components/buttons/SubmitButton';
import Editor, { Completion } from 'Components/Editor';
import { shouldSaveData } from 'Helpers/connection';
import { extractErrorMessage } from 'Helpers/utils';
import { useLoadAllSchemaEntities } from './graphql/loadAllSchemaEntities.generated';
import { useSQLShellContext } from '../SQLShellContext';
import { useRunQuery } from './graphql/runQuery.generated';

const minLines = 19;

const SQLEditor: React.FC = () => {
  const [value, setValue] = React.useState('');
  const { pushSnackbar } = useSnackbar();
  const {
    state: { selectedDatabase },
    dispatch,
  } = useSQLShellContext();

  const [runQuery, { loading: isSubmittingQueryRequest }] = useRunQuery({
    variables: {
      input: {
        databaseName: selectedDatabase,
        sql: value,
      },
    },
    onCompleted: data => {
      if (data.executeAsyncLogQuery.error) {
        dispatch({
          type: 'SET_ERROR',
          payload: { message: data.executeAsyncLogQuery.error.message },
        });
      } else {
        dispatch({ type: 'SET_QUERY_ID', payload: { queryId: data.executeAsyncLogQuery.queryId } });
      }
    },
    onError: error =>
      pushSnackbar({
        variant: 'error',
        title: "Couldn't execute your Query",
        description: extractErrorMessage(error),
      }),
  });

  // eslint-disable-next-line
  const [cancelQuery] = useRunQuery({
    variables: {
      input: {
        databaseName: selectedDatabase,
        sql: value,
      },
    },
    onCompleted: data => {
      if (data.executeAsyncLogQuery.error) {
        pushSnackbar({
          variant: 'error',
          title: "Couldn't cancel your Query",
          description: 'Your query will continue to be executed in the background',
        });
      } else {
        dispatch({ type: 'SET_QUERY_ID', payload: { queryId: null } });
      }
    },
    onError: () =>
      pushSnackbar({
        variant: 'error',
        title: "Couldn't cancel your Query",
        description: 'Your query will continue to be executed in the background',
      }),
  });

  // Fetch Autocomplete suggestions
  const { data: schemaData } = useLoadAllSchemaEntities({
    skip: shouldSaveData(),
    onError: error =>
      pushSnackbar({
        variant: 'warning',
        title: 'SQL autocomplete is disabled',
        description: extractErrorMessage(error),
      }),
  });

  React.useEffect(() => {
    if (isSubmittingQueryRequest) {
      dispatch({ type: 'RESET_ERROR' });
    }
  }, [isSubmittingQueryRequest]);

  // Create proper completion data
  const completions = React.useMemo(() => {
    const acc = new Set<Completion>();
    if (schemaData) {
      schemaData.listLogDatabases.forEach(database => {
        acc.add({ value: database.name, type: 'database' });
        database.tables.forEach(table => {
          acc.add({ value: table.name, type: 'table' });
          table.columns.forEach(column => {
            acc.add({ value: column.name, type: column.type });
          });
        });
      });
    }

    return [...acc];
  }, [schemaData]);

  return (
    <Box>
      <Editor
        fallback={<Box width="100%" bg="grey500" height={minLines * 19} />}
        placeholder="Run any SQL query. For example: select * from panther_logs.aws_alb;"
        minLines={minLines}
        mode="sql"
        width="100%"
        completions={completions}
        onChange={setValue}
        value={value}
      />
      <SubmitButton
        mt={6}
        disabled={!value || !selectedDatabase || isSubmittingQueryRequest}
        submitting={isSubmittingQueryRequest}
        onClick={() => runQuery()}
      >
        Run Query
      </SubmitButton>
    </Box>
  );
};

export default React.memo(SQLEditor);
