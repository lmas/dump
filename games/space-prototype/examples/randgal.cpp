// Original code: Bryan Mau
// www.cyberbrinedreams.com
// Modified by: Patrick Down
// www.codemoon.com

#include <windows.h>

#include <cmath>
#include <exception>
#include <string>
#include <vector>

using namespace std;

#define M_PI 3.1415927F

class Exception : public exception
{
public: 
    Exception( );
    Exception( const char* );
    Exception( const char*, const char*, long );
    virtual ~Exception();

    virtual const char *what() const throw();

protected:
    std::string m_sMessage;
    std::string m_sFile;
    long        m_nLine;
};


#define EXCEPTION( s ) Exception( s, __FILE__, __LINE__ )

struct Application_t
{
    Application_t() : m_hMainWindow( 0 ), m_hInst( 0 ), m_nStars( 10000 ), m_nArms( 0 ) {}
    
    HWND      m_hMainWindow;
    HINSTANCE m_hInst;
    string    m_sAppClassName;
    string    m_sAppName;

    int       m_nStars;
    int       m_nArms;
    int       m_nAngularSpread;

    vector<float> m_vfX;
    vector<float> m_vfY;
};

//////////////////////////////////////////////////////////////////////
// Private Forward References
//////////////////////////////////////////////////////////////////////
void ErrorBox( const char* cszError );
void ErrorBox( string& crsError );
void InitWindow( Application_t& rAppState );
void InitStars( Application_t& rAppState );
void Cleanup( Application_t& rAppState );
LRESULT CALLBACK WindowProc( HWND, UINT uMsg, WPARAM, LPARAM );

//////////////////////////////////////////////////////////////////////
// Public Variable
//////////////////////////////////////////////////////////////////////
//
// Name: g_AppState
//
// Description:
//   Global application state; can be referred to through global 
//   variable but it's preferred that it be accessed by explicitly 
//   passing it to the routines that need it.
//
// Special notes:
//
//////////////////////////////////////////////////////////////////////
Application_t g_AppState;


#include <time.h>
static time_t timStart = time( 0 );
static long   nFrame   = 0;

//////////////////////////////////////////////////////////////////////
// Public Function
//////////////////////////////////////////////////////////////////////
//
// Name: WinMain
//
// Description:
//   Main entry point to program
//
// Special notes:
//
//////////////////////////////////////////////////////////////////////
int WINAPI 
WinMain( 
        HINSTANCE hInstance, 
        HINSTANCE, 
        char*     szCmdLine, 
        int       nCmdShow 
        )
{
    int iRet = 0;

    try 
    {
        g_AppState.m_hInst         = hInstance;
        g_AppState.m_sAppClassName = "GalaxyGenApp";
        g_AppState.m_sAppName      = "Galaxy Generator";

        g_AppState.m_nStars         = 20000;
        g_AppState.m_nArms          = 6;
        g_AppState.m_nAngularSpread = 30;
 
        HWND hWnd = 0;

        InitWindow( g_AppState );
        InitStars( g_AppState );
 
        ShowWindow( g_AppState.m_hMainWindow, nCmdShow );

         // All set; get and process messages
        MSG msg;
        while ( GetMessage( &msg, NULL, 0, 0 ) )
        {
            if ( !TranslateAccelerator( hWnd, 0, &msg ) ) 
            {
                TranslateMessage( &msg );
                DispatchMessage( &msg );
            }
        }
        
        iRet = msg.wParam;
    }
    catch( exception* e )
    {
        string s = string( e->what() );
        ErrorBox( s ); 
    }
    catch(...)
    {
        ErrorBox( "Unspecified exception" ); 
    }

    Cleanup( g_AppState );

    {
        float fDelta = time( 0 ) - timStart;
        if ( fDelta != 0.0F )
        {
            char szMessage[100];
            sprintf( szMessage, "The cam was running at %f FPS", (nFrame/fDelta) );
//            ErrorBox( szMessage );
        }
    }

    return iRet;
}


//////////////////////////////////////////////////////////////////////
// Private Function
//////////////////////////////////////////////////////////////////////
//
// Name: InitWindow
//
// Description:
//   Helper function to create and initialize the main window
//
// Special notes:
//
//////////////////////////////////////////////////////////////////////
void InitWindow( 
    Application_t& rAppState
    )
{
    WNDCLASS wc;
    
    wc.lpszClassName = rAppState.m_sAppClassName.c_str();
    wc.hInstance     = rAppState.m_hInst;
    wc.lpfnWndProc   = WindowProc;
    wc.hCursor       = LoadCursor( 0, IDC_ARROW ) ;
    wc.hIcon         = 0;
    wc.lpszMenuName  = 0; // MAKEINTRESOURCE( IDR_MAIN_MENU );
    wc.hbrBackground = CreateSolidBrush( GetSysColor( COLOR_3DFACE ) );
    wc.style         = CS_HREDRAW | CS_VREDRAW ;
    wc.cbClsExtra    = 0 ;
    wc.cbWndExtra    = 0 ;

    if ( !RegisterClass( &wc ) ) 
    {
        throw new EXCEPTION( "Error registering window class for main window" );
    }

    string sCaption = rAppState.m_sAppName ;

    rAppState.m_hMainWindow = CreateWindowEx(
            0,
            rAppState.m_sAppClassName.c_str(),
            sCaption.c_str(),
            WS_CAPTION      |
            WS_SYSMENU      |
            WS_MINIMIZEBOX  |
            WS_MAXIMIZEBOX  |
            WS_THICKFRAME   |
            WS_CLIPCHILDREN |
            WS_OVERLAPPED,
            CW_USEDEFAULT, CW_USEDEFAULT,
//            CW_USEDEFAULT, CW_USEDEFAULT,
            800, 350,
//            800, 800,
            NULL,
            NULL,
            rAppState.m_hInst,
            0) ;

    if ( !rAppState.m_hMainWindow ) 
    {
        throw new EXCEPTION( "Error creating main application window" );
    }

}

float fRandom( float fMin, float fMax )
{
    static float fRandMax = 1.0f * RAND_MAX;

    float fRange = fMax - fMin;

    return fMin + fRange * (rand()/fRandMax);
}

float fLineRandom(float fRange)
{
  static float fRandMax = 1.0f * RAND_MAX;

  float fArea = fRange*fRange/2;
  float fP = fArea * (rand()/fRandMax);

  return fRange-sqrt(fRange*fRange - 2*fP);
}

float fHatRandom(float fRange)
{
  static float fRandMax = 1.0f * RAND_MAX;

  static float fArea = 4*atan(6.0);

  float fP = fArea * (rand()/fRandMax);

  return tan(fP/4)*fRange/6.0;
}

//////////////////////////////////////////////////////////////////////
// Private Function
//////////////////////////////////////////////////////////////////////
//
// Name: InitStars
//
// Description:
//   Helper function to generate a galaxy
//
// Special notes:
//
//////////////////////////////////////////////////////////////////////
void InitStars( 
    Application_t& rAppState
    )
{
    rAppState.m_vfX.clear();
    rAppState.m_vfY.clear();

    float fDeg2Rad       = M_PI / 180.0F;
    float fRadius        = 500.0F;

    float fArmAngle = (float)((360 / rAppState.m_nArms)%360);

    float fAngularSpread = 90/(rAppState.m_nArms);
//    float fAngularSpread = 200/(rAppState.m_nArms);

    for ( int i = 0; i < rAppState.m_nStars; i++ )
    {
//        float fR = fHatRandom(fRadius);
//        float fR = fLineRandom(fRadius);
//        float fQ = fLineRandom(fAngularSpread ) * (rand()&1 ? 1.0 : -1.0);
        float fR = fRandom(0.0, fRadius);
        float fQ = fRandom( 0.0, fAngularSpread ) * (rand()&1 ? 1.0 : -1.0);
        float fK = 1;

        float fA = (rand() % rAppState.m_nArms) * fArmAngle;

        float fX = fR * cos( fDeg2Rad * ( fA + fR * fK + fQ ) );
        float fY = fR * sin( fDeg2Rad * ( fA + fR * fK + fQ ) );

        rAppState.m_vfX.push_back( fX );
        rAppState.m_vfY.push_back( fY );
    }
}
//////////////////////////////////////////////////////////////////////
// Private Function
//////////////////////////////////////////////////////////////////////
//
// Name: Cleanup
//
// Description:
//   Cleanup the application when it exits
//
// Special notes:
//
//////////////////////////////////////////////////////////////////////
void Cleanup( 
    Application_t& rAppState
    )
{
    if ( rAppState.m_hMainWindow )
    {
        DestroyWindow( rAppState.m_hMainWindow );
    }
}

//////////////////////////////////////////////////////////////////////
// Private Functions
//////////////////////////////////////////////////////////////////////
//
// Name: ErrorBox
//
// Description:
//   Helper funtion to popup an error message box
//
// Special notes:
//
//////////////////////////////////////////////////////////////////////
void ErrorBox( 
              const char* cszError 
              )
{
    string s( cszError );

    ErrorBox( s );
}

void ErrorBox( string& crsError )
{
    string s = string( "MauCam experienced an error:\n\n" ) +
               crsError;

    MessageBox( GetDesktopWindow(), s.c_str(), "Error", MB_OK | MB_ICONSTOP );
}



//////////////////////////////////////////////////////////////////////
// Private Function
//////////////////////////////////////////////////////////////////////
//
// Name: WindowProc
//
// Description:
//   Window procedure for main window
//
// Special notes:
//
//////////////////////////////////////////////////////////////////////
LRESULT CALLBACK WindowProc(
  HWND   hWnd,    // handle to window
  UINT   uMsg,    // message identifier
  WPARAM wParam,  // first message parameter
  LPARAM lParam   // second message parameter
)
{
    switch ( uMsg )
    {
//    case WM_CLOSE:
//        return 0;


#if 0
    case WM_COMMAND:
        {
            switch ( LOWORD( wParam ) )
            {
            case IDM_FILE_NEW:
                PopupNewView( g_AppState );
                return 0;

            case IDM_FILE_EXIT:
                SendMessage( hWnd, WM_CLOSE, 0, 0 );
                return 0;

            case IDM_VIEW_FTP_SETTINGS:
                //
                //  If there are no ftp sessions create one
                //
                if ( 0 == g_AppState.m_vFtpState.size() )
                {
                    Ftp_t* pFtp = new Ftp_t;
                    g_AppState.m_vFtpState.push_back( pFtp );
                }
                PopupFtpSettings( g_AppState, *(g_AppState.m_vFtpState[0]) );
                WriteFtpProfile( g_AppState );
                return 0;

            case IDM_VIEW_CAMERASETTINGS:
                PopupCamSettings( g_AppState );
                return 0;

            case IDM_HELP_ABOUT:
                PopupAbout( g_AppState );
                return 0;
            }
        }
        break;

    case WM_TIMER:
        {
#if 1
#if 0
            nFrame++;
            
            std::vector<VidCap*>::iterator it = 0;
            for ( it = g_AppState.m_vpCam.begin(); 
                  it != g_AppState.m_vpCam.end(); 
                  ++it )
            {
                DWORD dwId = 0;
                HANDLE hThread = CreateThread( 0, 0, GrabThreadFunc, (*it), 0, &dwId );
            }
#else
            nFrame++;
            
            std::vector<VidFeed*>::iterator it = 0;
            for ( it = g_AppState.m_vpFeed.begin(); 
                  it != g_AppState.m_vpFeed.end(); 
                  ++it )
            {
                (*it)->Grab();
            }
#endif
#endif
        }
        return 0;
#endif

    case WM_PAINT:
        {
            PAINTSTRUCT ps;
            HDC hdc = BeginPaint( hWnd, &ps );
            {
                RECT rcClient;
                GetClientRect( hWnd, &rcClient );

                FillRect( hdc, &rcClient, (HBRUSH)GetStockObject( BLACK_BRUSH ) );

                int nOrgX = ( rcClient.right  - rcClient.left ) / 2;
                int nOrgY = ( rcClient.bottom - rcClient.top ) / 2;

                for ( int i = 0; i < g_AppState.m_nStars; i++ )
                {
                    SetPixel( 
                        hdc,
                        nOrgX + (int)g_AppState.m_vfX[i],
                        nOrgY + (int)g_AppState.m_vfY[i],
                        RGB( 128, 128, 128 )
                        );
                }

            }
            EndPaint( hWnd, &ps );
        }
        return 0;

    case WM_SIZE:
        {
            RECT rcClient;
            GetClientRect( g_AppState.m_hMainWindow, &rcClient );
                   }
        return 0;

    case WM_DESTROY:
        g_AppState.m_hMainWindow = 0; // keep window from being destroyed twice on cleanup
        PostQuitMessage( 0 );
        return 0;
    }

    return DefWindowProc( hWnd, uMsg, wParam, lParam );
}


Exception::Exception() : m_nLine( -1 )
{
}

Exception::Exception( 
   const char* cszMessage
   ) : m_sMessage( cszMessage ),
       m_nLine( -1 )
{
}

Exception::Exception( 
   const char* cszMessage, 
   const char* cszFile, 
   long        nLine 
   ) : m_sMessage( cszMessage ),
       m_sFile ( cszFile ),
       m_nLine( nLine )
{
}

const char*
Exception::what() const
{
    static string sMessage( "Exception:\n" );

    if ( m_sMessage.size() > 0 )
    {
        sMessage += "\n  Error: ";
        sMessage += m_sMessage;
        sMessage += "\n";
    }


    if ( m_nLine > 0 )
    {
        char szProgrammerInfo[1024];
        sprintf( szProgrammerInfo, "\n  Programmer info: %s(%i)\n", 
            m_sFile.c_str(), m_nLine );
        sMessage += szProgrammerInfo;
    }

    return sMessage.c_str();
}

Exception::~Exception()
{
}
