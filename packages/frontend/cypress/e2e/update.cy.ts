describe('Update an animal', () => {
  beforeEach(() => {
    cy.intercept('GET', 'http://localhost:3000/animals', { fixture: 'animals.json' }).as('getAnimals')
    cy.visit('/')
    cy.wait('@getAnimals')
  })

  it('update an animal', () => {
    cy.get('table tbody tr')
      .contains('Leo')
      .parents('tr')
      .within(() => {
        cy.get('svg').click({ force: true })
      })

    cy.get('.modal').should('be.visible')

    cy.get('input[id="name"]').should('have.value', 'Leo')

    cy.get('input[id="name"]').clear().type('Leonardo')

    cy.get('button').contains('Update').click()

    cy.get('.modal').should('not.exist')

    // cy.get('table').contains('Leonardo').should('exist')
    // cy.get('table').contains('Leo').should('not.exist')
  })
})
