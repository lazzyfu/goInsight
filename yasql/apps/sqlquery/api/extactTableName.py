# Create your tests here.


# code from https://github.com/macbre/sql-metadata
import re
from typing import List

import sqlparse
from sqlparse.sql import TokenList
from sqlparse.tokens import Name, Whitespace


def unique(_list: List) -> List:
    """
    Makes the list have unique items only and maintains the order
    list(set()) won't provide that
    :type _list list
    :rtype: list
    """
    ret = []

    for item in _list:
        if item not in ret:
            ret.append(item)

    return ret


def preprocess_query(query: str) -> str:
    """
    Perform initial query cleanup
    :type query str
    :rtype str
    """
    # 0. remove newlines
    query = query.replace('\n', ' ')

    # 1. remove aliases
    # FROM `dimension_wikis` `dw`
    # INNER JOIN `fact_wam_scores` `fwN`
    query = re.sub(r'(\s(FROM|JOIN)\s`[^`]+`)\s`[^`]+`', r'\1', query, flags=re.IGNORECASE)

    # 2. `database`.`table` notation -> database.table
    query = re.sub(r'`([^`]+)`\.`([^`]+)`', r'\1.\2', query)

    # 2. database.table notation -> table
    # query = re.sub(r'([a-z_0-9]+)\.([a-z_0-9]+)', r'\2', query, flags=re.IGNORECASE)

    return query


def get_query_tokens(query: str) -> List[sqlparse.sql.Token]:
    """
    :type query str
    :rtype: list[sqlparse.sql.Token]
    """
    query = preprocess_query(query)
    parsed = sqlparse.parse(query)

    # handle empty queries (#12)
    if not parsed:
        return []

    tokens = TokenList(parsed[0].tokens).flatten()
    # print([(token.value, token.ttype) for token in tokens])

    return [token for token in tokens if token.ttype is not Whitespace]


def _update_table_names(tables: List[str],
                        tokens: List[sqlparse.sql.Token],
                        index: int,
                        last_keyword: str) -> List[str]:
    """
    Return new table names matching database.table or database.schema.table notation
    :type tables list[str]
    :type tokens list[sqlparse.sql.Token]
    :type index int
    :type last_keyword str
    :rtype: list[str]
    """

    token = tokens[index]
    last_token = tokens[index - 1].value.upper() if index > 0 else None
    next_token = tokens[index + 1].value.upper() if index + 1 < len(tokens) else None

    if last_keyword in ['FROM', 'JOIN', 'INNER JOIN', 'FULL JOIN', 'FULL OUTER JOIN',
                        'LEFT JOIN', 'RIGHT JOIN', 'STRAIGHT_JOIN',
                        'LEFT OUTER JOIN', 'RIGHT OUTER JOIN'] \
            and last_token not in ['AS'] \
            and token.value not in ['AS', 'SELECT']:
        if last_token == '.' and next_token != '.':
            # we have database.table notation example
            table_name = '{}.{}'.format(tokens[index - 2], tokens[index])
            if len(tables) > 0:
                tables[-1] = table_name
            else:
                tables.append(table_name)

        schema_notation_match = (Name, '.', Name, '.', Name)
        schema_notation_tokens = (tokens[index - 4].ttype,
                                  tokens[index - 3].value,
                                  tokens[index - 2].ttype,
                                  tokens[index - 1].value,
                                  tokens[index].ttype) if len(tokens) > 4 else None
        if schema_notation_tokens == schema_notation_match:
            # we have database.schema.table notation example
            table_name = '{}.{}.{}'.format(
                tokens[index - 4], tokens[index - 2], tokens[index])
            if len(tables) > 0:
                tables[-1] = table_name
            else:
                tables.append(table_name)
        elif tokens[index - 1].value.upper() not in [',', last_keyword]:
            # it's not a list of tables, e.g. SELECT * FROM foo, bar
            # hence, it can be the case of alias without AS, e.g. SELECT * FROM foo bar
            pass
        else:
            table_name = str(token.value.strip('`'))
            tables.append(table_name)

    return tables


def get_query_tables(query: str) -> List[str]:
    """
    :type query str
    :rtype: list[str]
    """
    tables = []
    last_keyword = None

    table_syntax_keywords = [
        # SELECT queries
        'FROM', 'WHERE', 'JOIN', 'INNER JOIN', 'FULL JOIN', 'FULL OUTER JOIN',
        'LEFT OUTER JOIN', 'RIGHT OUTER JOIN', 'STRAIGHT_JOIN',
        'LEFT JOIN', 'RIGHT JOIN', 'ON'
    ]

    # print(query, get_query_tokens(query))
    query = query.replace('"', '')
    tokens = get_query_tokens(query)

    for index, token in enumerate(tokens):
        # print([token, token.ttype, last_token, last_keyword])
        if token.is_keyword and token.value.upper() in table_syntax_keywords:
            # keep the name of the last keyword, the next one can be a table name
            last_keyword = token.value.upper()
            # print('keyword', last_keyword)
        elif str(token) == '(':
            # reset the last_keyword for INSERT `foo` VALUES(id, bar) ...
            last_keyword = None
        elif token.is_keyword and str(token) in ['FORCE', 'ORDER', 'GROUP BY']:
            # reset the last_keyword for queries like:
            # "SELECT x FORCE INDEX"
            # "SELECT x ORDER BY"
            # "SELECT x FROM y GROUP BY x"
            last_keyword = None
        elif token.is_keyword and str(token) == 'SELECT' and last_keyword in ['INTO', 'TABLE']:
            # reset the last_keyword for "INSERT INTO SELECT" and "INSERT TABLE SELECT" queries
            last_keyword = None
        elif token.ttype is Name or token.is_keyword:
            tables = _update_table_names(tables, tokens, index, last_keyword)

    return unique(tables)
